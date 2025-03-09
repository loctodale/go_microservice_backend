namespace Order.API.Repository;

public interface IOrderRepository
{
    Task<List<Entities.Order>> GetOrders();
    Task<Entities.Order> GetOrder(int id);
}