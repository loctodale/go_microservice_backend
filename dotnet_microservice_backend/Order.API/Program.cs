using Order.API.Config;
using Order.API.Services;

var builder = WebApplication.CreateBuilder(args);
var mapper = OrderAutoMapper.InitAutoMapper();

// Add services to the container.
builder.Services.AddGrpc(optiion =>
{
    optiion.EnableDetailedErrors = true;
}).AddJsonTranscoding();
var app = builder.Build();

// Configure the HTTP request pipeline.
app.MapGrpcService<GreeterService>();
app.MapGrpcService<OrderService>();
app.MapGet("/",
    () =>
        "Communication with gRPC endpoints must be made through a gRPC client. To learn how to create a client, visit: https://go.microsoft.com/fwlink/?linkid=2086909");

app.Run();