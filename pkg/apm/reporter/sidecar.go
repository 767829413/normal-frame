package reporter

import (
	"bytes"
	"encoding/json"
	"log"
	"net"

	"github.com/767829413/normal-frame/fork/SkyAPM/go2sky"
	agentv3 "github.com/767829413/normal-frame/fork/SkyAPM/go2sky/reporter/grpc/language-agent"
)

// NewSidecarReporter create a new reporter to send data to sidecar. Only one backend address is allowed.
func NewSidecarReporter(serverAddr string) (go2sky.Reporter, error) {
	r := &sidecarReporter{}
	conn, err := net.Dial("unix", serverAddr)
	if err != nil {
		return nil, err
	}
	r.conn = conn
	return r, nil
}

type sidecarReporter struct {
	service         string
	serviceInstance string
	instanceProps   map[string]string
	conn            net.Conn
	bootFlag        bool
}

func (r *sidecarReporter) Boot(service string, serviceInstance string) {
	r.service = service
	r.serviceInstance = serviceInstance
	r.bootFlag = true
}

func (r *sidecarReporter) Send(spans []go2sky.ReportedSpan) {
	spanSize := len(spans)
	if spanSize < 1 {
		return
	}
	rootSpan := spans[spanSize-1]
	rootCtx := rootSpan.Context()
	segmentObject := &agentv3.SegmentObject{
		TraceId:         rootCtx.TraceID,
		TraceSegmentId:  rootCtx.SegmentID,
		Spans:           make([]*agentv3.SpanObject, spanSize),
		Service:         r.service,
		ServiceInstance: r.serviceInstance,
	}
	for i, s := range spans {
		spanCtx := s.Context()
		segmentObject.Spans[i] = &agentv3.SpanObject{
			SpanId:        spanCtx.SpanID,
			ParentSpanId:  spanCtx.ParentSpanID,
			StartTime:     s.StartTime(),
			EndTime:       s.EndTime(),
			OperationName: s.OperationName(),
			Peer:          s.Peer(),
			SpanType:      s.SpanType(),
			SpanLayer:     s.SpanLayer(),
			ComponentId:   s.ComponentID(),
			IsError:       s.IsError(),
			Tags:          s.Tags(),
			Logs:          s.Logs(),
		}
		srr := make([]*agentv3.SegmentReference, 0)
		if i == (spanSize-1) && spanCtx.ParentSpanID > -1 {
			srr = append(srr, &agentv3.SegmentReference{
				RefType:               agentv3.RefType_CrossThread,
				TraceId:               spanCtx.TraceID,
				ParentTraceSegmentId:  spanCtx.ParentSegmentID,
				ParentSpanId:          spanCtx.ParentSpanID,
				ParentService:         r.service,
				ParentServiceInstance: r.serviceInstance,
			})
		}
		if len(s.Refs()) > 0 {
			for _, tc := range s.Refs() {
				srr = append(srr, &agentv3.SegmentReference{
					RefType:                  agentv3.RefType_CrossProcess,
					TraceId:                  spanCtx.TraceID,
					ParentTraceSegmentId:     tc.ParentSegmentID,
					ParentSpanId:             tc.ParentSpanID,
					ParentService:            tc.ParentService,
					ParentServiceInstance:    tc.ParentServiceInstance,
					ParentEndpoint:           tc.ParentEndpoint,
					NetworkAddressUsedAtPeer: tc.AddressUsedAtClient,
				})
			}
		}
		segmentObject.Spans[i].Refs = srr
	}
	defer func() {
		// recover the panic caused by close sendCh
		if err := recover(); err != nil {
			log.Printf("reporter segment err %v", err)
		}
	}()
	// 组合成符合sidecar格式的数据
	var segment = &SegmentObject{
		TraceId: segmentObject.TraceId,
		Segment: Segment{
			TraceSegmentId: segmentObject.TraceSegmentId,
			IsSizeLimited:  segmentObject.IsSizeLimited,
			Spans:          nil,
		},
		Service:         segmentObject.Service,
		ServiceInstance: segmentObject.ServiceInstance,
	}
	for _, value := range segmentObject.Spans {
		span := &SpanObject{
			SpanId:        value.SpanId,
			ParentSpanId:  value.ParentSpanId,
			StartTime:     value.StartTime,
			EndTime:       value.EndTime,
			OperationName: value.OperationName,
			Peer:          value.Peer,
			SpanType:      value.SpanType,
			SpanLayer:     value.SpanLayer,
			ComponentId:   value.ComponentId,
			IsError:       value.IsError,
			//Tags:          nil,
			SkipAnalysis: value.SkipAnalysis,
		}
		tags := make(map[string]string)
		for _, v := range value.Tags {
			tags[v.Key] = v.Value
		}
		span.Tags = tags
		segment.Segment.Spans = append(segment.Segment.Spans, span)
	}
	byteSli, _ := json.Marshal(segment)
	var buf bytes.Buffer
	buf.WriteByte(49)
	buf.Write(byteSli)
	buf.WriteByte(10)
	if _, err := r.conn.Write(buf.Bytes()); err != nil {
		log.Printf("conn Write error %v", err)
		r.Close()
	}
}

type SegmentObject struct {
	TraceId         string  `json:"traceId,omitempty"`
	Segment         Segment `json:"segment,omitempty"`
	Service         string  `json:"service,omitempty"`
	ServiceInstance string  `json:"serviceInstance,omitempty"`
}

type Segment struct {
	TraceSegmentId string        `json:"traceSegmentId,omitempty"`
	IsSizeLimited  bool          `json:"isSizeLimited,omitempty"`
	Spans          []*SpanObject `json:"spans,omitempty"`
}

type SpanObject struct {
	SpanId        int32             `protobuf:"varint,1,opt,name=spanId,proto3" json:"spanId,omitempty"`
	ParentSpanId  int32             `protobuf:"varint,2,opt,name=parentSpanId,proto3" json:"parentSpanId,omitempty"`
	StartTime     int64             `protobuf:"varint,3,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime       int64             `protobuf:"varint,4,opt,name=endTime,proto3" json:"endTime,omitempty"`
	OperationName string            `protobuf:"bytes,6,opt,name=operationName,proto3" json:"operationName,omitempty"`
	Peer          string            `protobuf:"bytes,7,opt,name=peer,proto3" json:"peer,omitempty"`
	SpanType      agentv3.SpanType  `protobuf:"varint,8,opt,name=spanType,proto3,enum=SpanType" json:"spanType,omitempty"`
	SpanLayer     agentv3.SpanLayer `protobuf:"varint,9,opt,name=spanLayer,proto3,enum=SpanLayer" json:"spanLayer,omitempty"`
	ComponentId   int32             `protobuf:"varint,10,opt,name=componentId,proto3" json:"componentId,omitempty"`
	IsError       bool              `protobuf:"varint,11,opt,name=isError,proto3" json:"isError,omitempty"`
	Tags          map[string]string `protobuf:"bytes,12,rep,name=tags,proto3" json:"tags,omitempty"`
	SkipAnalysis  bool              `protobuf:"varint,14,opt,name=skipAnalysis,proto3" json:"skipAnalysis,omitempty"`
}

func (r *sidecarReporter) Close() {
	r.closeConn()
}

func (r *sidecarReporter) closeConn() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Println(err)
		}
	}
}
