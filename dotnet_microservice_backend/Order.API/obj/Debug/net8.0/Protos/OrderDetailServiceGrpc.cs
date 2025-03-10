// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: Protos/order_detail_service.proto
// </auto-generated>
#pragma warning disable 0414, 1591, 8981, 0612
#region Designer generated code

using grpc = global::Grpc.Core;

namespace order_detail_service {
  public static partial class OrderDetailService
  {
    static readonly string __ServiceName = "orderdetailservice.OrderDetailService";

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static void __Helper_SerializeMessage(global::Google.Protobuf.IMessage message, grpc::SerializationContext context)
    {
      #if !GRPC_DISABLE_PROTOBUF_BUFFER_SERIALIZATION
      if (message is global::Google.Protobuf.IBufferMessage)
      {
        context.SetPayloadLength(message.CalculateSize());
        global::Google.Protobuf.MessageExtensions.WriteTo(message, context.GetBufferWriter());
        context.Complete();
        return;
      }
      #endif
      context.Complete(global::Google.Protobuf.MessageExtensions.ToByteArray(message));
    }

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static class __Helper_MessageCache<T>
    {
      public static readonly bool IsBufferMessage = global::System.Reflection.IntrospectionExtensions.GetTypeInfo(typeof(global::Google.Protobuf.IBufferMessage)).IsAssignableFrom(typeof(T));
    }

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static T __Helper_DeserializeMessage<T>(grpc::DeserializationContext context, global::Google.Protobuf.MessageParser<T> parser) where T : global::Google.Protobuf.IMessage<T>
    {
      #if !GRPC_DISABLE_PROTOBUF_BUFFER_SERIALIZATION
      if (__Helper_MessageCache<T>.IsBufferMessage)
      {
        return parser.ParseFrom(context.PayloadAsReadOnlySequence());
      }
      #endif
      return parser.ParseFrom(context.PayloadAsNewBuffer());
    }

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::order_detail_service.RequestGetOrderDetail> __Marshaller_orderdetailservice_RequestGetOrderDetail = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::order_detail_service.RequestGetOrderDetail.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::order_detail_service.ResponseGetOrderDetail> __Marshaller_orderdetailservice_ResponseGetOrderDetail = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::order_detail_service.ResponseGetOrderDetail.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::order_detail_service.RequestCreateOrderDetail> __Marshaller_orderdetailservice_RequestCreateOrderDetail = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::order_detail_service.RequestCreateOrderDetail.Parser));
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Marshaller<global::order_detail_service.EmptyResponse> __Marshaller_orderdetailservice_EmptyResponse = grpc::Marshallers.Create(__Helper_SerializeMessage, context => __Helper_DeserializeMessage(context, global::order_detail_service.EmptyResponse.Parser));

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Method<global::order_detail_service.RequestGetOrderDetail, global::order_detail_service.ResponseGetOrderDetail> __Method_GetOrderDetail = new grpc::Method<global::order_detail_service.RequestGetOrderDetail, global::order_detail_service.ResponseGetOrderDetail>(
        grpc::MethodType.Unary,
        __ServiceName,
        "GetOrderDetail",
        __Marshaller_orderdetailservice_RequestGetOrderDetail,
        __Marshaller_orderdetailservice_ResponseGetOrderDetail);

    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    static readonly grpc::Method<global::order_detail_service.RequestCreateOrderDetail, global::order_detail_service.EmptyResponse> __Method_CreateOrderDetail = new grpc::Method<global::order_detail_service.RequestCreateOrderDetail, global::order_detail_service.EmptyResponse>(
        grpc::MethodType.Unary,
        __ServiceName,
        "CreateOrderDetail",
        __Marshaller_orderdetailservice_RequestCreateOrderDetail,
        __Marshaller_orderdetailservice_EmptyResponse);

    /// <summary>Service descriptor</summary>
    public static global::Google.Protobuf.Reflection.ServiceDescriptor Descriptor
    {
      get { return global::order_detail_service.OrderDetailServiceReflection.Descriptor.Services[0]; }
    }

    /// <summary>Base class for server-side implementations of OrderDetailService</summary>
    [grpc::BindServiceMethod(typeof(OrderDetailService), "BindService")]
    public abstract partial class OrderDetailServiceBase
    {
      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::System.Threading.Tasks.Task<global::order_detail_service.ResponseGetOrderDetail> GetOrderDetail(global::order_detail_service.RequestGetOrderDetail request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

      [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
      public virtual global::System.Threading.Tasks.Task<global::order_detail_service.EmptyResponse> CreateOrderDetail(global::order_detail_service.RequestCreateOrderDetail request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

    }

    /// <summary>Creates service definition that can be registered with a server</summary>
    /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public static grpc::ServerServiceDefinition BindService(OrderDetailServiceBase serviceImpl)
    {
      return grpc::ServerServiceDefinition.CreateBuilder()
          .AddMethod(__Method_GetOrderDetail, serviceImpl.GetOrderDetail)
          .AddMethod(__Method_CreateOrderDetail, serviceImpl.CreateOrderDetail).Build();
    }

    /// <summary>Register service method with a service binder with or without implementation. Useful when customizing the service binding logic.
    /// Note: this method is part of an experimental API that can change or be removed without any prior notice.</summary>
    /// <param name="serviceBinder">Service methods will be bound by calling <c>AddMethod</c> on this object.</param>
    /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
    [global::System.CodeDom.Compiler.GeneratedCode("grpc_csharp_plugin", null)]
    public static void BindService(grpc::ServiceBinderBase serviceBinder, OrderDetailServiceBase serviceImpl)
    {
      serviceBinder.AddMethod(__Method_GetOrderDetail, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::order_detail_service.RequestGetOrderDetail, global::order_detail_service.ResponseGetOrderDetail>(serviceImpl.GetOrderDetail));
      serviceBinder.AddMethod(__Method_CreateOrderDetail, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::order_detail_service.RequestCreateOrderDetail, global::order_detail_service.EmptyResponse>(serviceImpl.CreateOrderDetail));
    }

  }
}
#endregion
