using AutoMapper;

namespace Order.API.Config;

public class OrderAutoMapper
{
  public static IMapper InitAutoMapper()
  {
    var config = new MapperConfiguration(cfg =>
    {
      OrderConfig.CreateMap(cfg);
    }); 
    return config.CreateMapper();
  }
  public class OrderConfig()
  {
    public static void CreateMap(IMapperConfigurationExpression cfg)
    {
      cfg.CreateMap<Entities.Order, order_service.Order>().ReverseMap();
      // cfg.CreateMap<List<Entities.Order>, List<order_service.Order>>().ReverseMap();
    }
  }
}