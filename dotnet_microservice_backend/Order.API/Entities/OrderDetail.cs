using System;
using System.Collections.Generic;

namespace Order.API.Entities;

public partial class OrderDetail
{
    public int OrderDetailId { get; set; }

    public int OrderId { get; set; }

    public int ProductId { get; set; }

    public int Quantity { get; set; }

    public int PriceEachItem { get; set; }

    public int TotalPrice { get; set; }

    public DateOnly? CreatedDate { get; set; }

    public DateOnly? UpdatedDate { get; set; }

    public DateOnly? DeteledDate { get; set; }

    public virtual Order Order { get; set; } = null!;
}
