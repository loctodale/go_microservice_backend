# Build stage
FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:8.0-alpine AS build
ARG TARGETARCH
WORKDIR /source

# Install compatibility libraries and protoc
RUN apk add --no-cache libc6-compat protobuf

# Copy project file and restore
COPY --link Order.API/*.csproj .
RUN dotnet restore -a $TARGETARCH

# Copy source code and publish
COPY --link Order.API/. .
RUN dotnet publish --no-restore -a $TARGETARCH -o /app --no-cache /restore

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine
WORKDIR /app
EXPOSE 8080

# Copy published app from build stage
COPY --link --from=build /app .

# Run the app
ENTRYPOINT ["dotnet", "Order.API.dll"]
