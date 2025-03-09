using AutoMapper;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using order_service;
using Order.API.Repository;
using Empty = order_service.Empty;

namespace Order.API.Services;

public class OrderService : order_service.OrderService.OrderServiceBase
{
    private readonly ILogger<OrderService> _logger;
    private readonly IOrderRepository _orderRepository;
    private readonly IMapper _mapper;

    public OrderService(ILogger<OrderService> logger, IOrderRepository orderRepository, IMapper mapper)
    {
        _logger = logger;
        _orderRepository = orderRepository;
        _mapper = mapper;
    }

    public override async Task<ResponseGetAllOrder> GetAllOrder(Empty request, ServerCallContext context)
    {
        List<Entities.Order> result = await _orderRepository.GetOrders();

        var data = result.Select(x => _mapper.Map<Entities.Order, order_service.Order>(x)).ToList();
        return await Task.FromResult(new ResponseGetAllOrder
        {
            Message = "Success",
            Data = {data},
            StatusCode = 200
        });
    }

    public override async Task<ResponseGetOrderById> GetOrderById(RequestGetOrderById request, ServerCallContext context)
    {
        Entities.Order result = await _orderRepository.GetOrder(request.OrderId);
        return await Task.FromResult(new ResponseGetOrderById
        {
            Message = "Success",
            StatusCode = 200,
            Data = _mapper.Map<Entities.Order, order_service.Order>(result)
        });
    }
}