namespace Order.API.Repository.OrderDetail;

public interface IOrderDetailRepository
{
    Task<Entities.OrderDetail> GetOrderDetail(int orderDetailId);
}