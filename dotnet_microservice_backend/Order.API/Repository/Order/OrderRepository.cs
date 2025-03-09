using Microsoft.EntityFrameworkCore;
using Order.API.Context;

namespace Order.API.Repository;

public class OrderRepository : IOrderRepository
{
    private OrderDbContext _context;
    
    public OrderRepository(OrderDbContext context) => _context = context;


    public async Task<List<Entities.Order>> GetOrders()
    {
        return await _context.Orders.ToListAsync();
    }

    public async Task<Entities.Order> GetOrder(int id)
    {
        return await _context.Orders.FindAsync(id);
    }
}