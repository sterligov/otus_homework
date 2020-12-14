// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	GetEventByID(ctx context.Context, in *EventID, opts ...grpc.CallOption) (*Event, error)
	CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Inserted, error)
	UpdateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Affected, error)
	DeleteEvent(ctx context.Context, in *EventID, opts ...grpc.CallOption) (*Affected, error)
	GetUserDayEvents(ctx context.Context, in *UserPeriodEventRequest, opts ...grpc.CallOption) (*Events, error)
	GetUserWeekEvents(ctx context.Context, in *UserPeriodEventRequest, opts ...grpc.CallOption) (*Events, error)
	GetUserMonthEvents(ctx context.Context, in *UserPeriodEventRequest, opts ...grpc.CallOption) (*Events, error)
	Health(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HealthResponse, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) GetEventByID(ctx context.Context, in *EventID, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/event.EventService/GetEventByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) CreateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Inserted, error) {
	out := new(Inserted)
	err := c.cc.Invoke(ctx, "/event.EventService/CreateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) UpdateEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*Affected, error) {
	out := new(Affected)
	err := c.cc.Invoke(ctx, "/event.EventService/UpdateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) DeleteEvent(ctx context.Context, in *EventID, opts ...grpc.CallOption) (*Affected, error) {
	out := new(Affected)
	err := c.cc.Invoke(ctx, "/event.EventService/DeleteEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetUserDayEvents(ctx context.Context, in *UserPeriodEventRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, "/event.EventService/GetUserDayEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetUserWeekEvents(ctx context.Context, in *UserPeriodEventRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, "/event.EventService/GetUserWeekEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) GetUserMonthEvents(ctx context.Context, in *UserPeriodEventRequest, opts ...grpc.CallOption) (*Events, error) {
	out := new(Events)
	err := c.cc.Invoke(ctx, "/event.EventService/GetUserMonthEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventServiceClient) Health(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, "/event.EventService/Health", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility
type EventServiceServer interface {
	GetEventByID(context.Context, *EventID) (*Event, error)
	CreateEvent(context.Context, *Event) (*Inserted, error)
	UpdateEvent(context.Context, *Event) (*Affected, error)
	DeleteEvent(context.Context, *EventID) (*Affected, error)
	GetUserDayEvents(context.Context, *UserPeriodEventRequest) (*Events, error)
	GetUserWeekEvents(context.Context, *UserPeriodEventRequest) (*Events, error)
	GetUserMonthEvents(context.Context, *UserPeriodEventRequest) (*Events, error)
	Health(context.Context, *Empty) (*HealthResponse, error)
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (UnimplementedEventServiceServer) GetEventByID(context.Context, *EventID) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventByID not implemented")
}
func (UnimplementedEventServiceServer) CreateEvent(context.Context, *Event) (*Inserted, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedEventServiceServer) UpdateEvent(context.Context, *Event) (*Affected, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedEventServiceServer) DeleteEvent(context.Context, *EventID) (*Affected, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (UnimplementedEventServiceServer) GetUserDayEvents(context.Context, *UserPeriodEventRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserDayEvents not implemented")
}
func (UnimplementedEventServiceServer) GetUserWeekEvents(context.Context, *UserPeriodEventRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserWeekEvents not implemented")
}
func (UnimplementedEventServiceServer) GetUserMonthEvents(context.Context, *UserPeriodEventRequest) (*Events, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserMonthEvents not implemented")
}
func (UnimplementedEventServiceServer) Health(context.Context, *Empty) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	s.RegisterService(&_EventService_serviceDesc, srv)
}

func _EventService_GetEventByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetEventByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetEventByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetEventByID(ctx, req.(*EventID))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/CreateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).CreateEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/UpdateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).UpdateEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/DeleteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).DeleteEvent(ctx, req.(*EventID))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetUserDayEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserPeriodEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetUserDayEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetUserDayEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetUserDayEvents(ctx, req.(*UserPeriodEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetUserWeekEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserPeriodEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetUserWeekEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetUserWeekEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetUserWeekEvents(ctx, req.(*UserPeriodEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_GetUserMonthEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserPeriodEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).GetUserMonthEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/GetUserMonthEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).GetUserMonthEvents(ctx, req.(*UserPeriodEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EventService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.EventService/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).Health(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "event.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEventByID",
			Handler:    _EventService_GetEventByID_Handler,
		},
		{
			MethodName: "CreateEvent",
			Handler:    _EventService_CreateEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _EventService_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _EventService_DeleteEvent_Handler,
		},
		{
			MethodName: "GetUserDayEvents",
			Handler:    _EventService_GetUserDayEvents_Handler,
		},
		{
			MethodName: "GetUserWeekEvents",
			Handler:    _EventService_GetUserWeekEvents_Handler,
		},
		{
			MethodName: "GetUserMonthEvents",
			Handler:    _EventService_GetUserMonthEvents_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _EventService_Health_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/event_service.proto",
}
