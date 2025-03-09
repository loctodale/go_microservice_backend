using System;
using System.Collections.Generic;

namespace Order.API.Entities;

public partial class Order
{
    public int OrderId { get; set; }

    public int OrderUserId { get; set; }

    public int? OrderShippingId { get; set; }

    public string? OrderPayment { get; set; }

    public string? OrderTrackingNumber { get; set; }

    public string? OrderStatus { get; set; }

    public DateOnly? CreatedDate { get; set; }

    public DateOnly? UpdatedDate { get; set; }

    public DateOnly? DeteledDate { get; set; }

    public virtual ICollection<OrderDetail> OrderDetails { get; set; } = new List<OrderDetail>();
}
