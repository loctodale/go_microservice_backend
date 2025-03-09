using AutoMapper;
using Order.API.Context;

namespace Order.API.Repository.OrderDetail;

public class OrderDetailRepository : IOrderDetailRepository
{
    private readonly OrderDbContext _context;
    
    public OrderDetailRepository(OrderDbContext context) => _context = context;

    public async Task<Entities.OrderDetail> GetOrderDetail(int orderDetailId)
    {
        return await _context.OrderDetails.FindAsync(orderDetailId);
    }
}