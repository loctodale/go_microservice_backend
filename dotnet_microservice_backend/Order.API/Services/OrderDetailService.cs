using AutoMapper;
using Grpc.Core;
using order_detail_service;
using Order.API.Repository.OrderDetail;

namespace Order.API.Services;

public class OrderDetailService  : order_detail_service.OrderDetailService.OrderDetailServiceBase
{
    private readonly IOrderDetailRepository _orderDetailRepository;
    private readonly IMapper _mapper;
    private readonly ILogger<OrderDetailService> _logger;

    public OrderDetailService(IOrderDetailRepository orderDetailRepository, IMapper mapper,
        ILogger<OrderDetailService> logger)
    {
        _logger = logger;
        _orderDetailRepository = orderDetailRepository;
        _mapper = mapper;
    }

    public override async Task<ResponseGetOrderDetail> GetOrderDetail(RequestGetOrderDetail request, ServerCallContext context)
    {
        Entities.OrderDetail result = await _orderDetailRepository.GetOrderDetail(request.OrderDetailId);
        return await Task.FromResult(new ResponseGetOrderDetail
        {
            StatusCode = StatusCodes.Status200OK,
            Message = "OK",
            Data = _mapper.Map<OrderDetail>(result)
        });
    }
}