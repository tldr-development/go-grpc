//
// DO NOT EDIT.
// swift-format-ignore-file
//
// Generated by the protocol buffer compiler.
// Source: account.proto
//
import GRPC
import NIO
import NIOConcurrencyHelpers
import SwiftProtobuf


/// Usage: instantiate `Account_AddServiceClient`, then call methods of this protocol to make API calls.
internal protocol Account_AddServiceClientProtocol: GRPCClient {
  var serviceName: String { get }
  var interceptors: Account_AddServiceClientInterceptorFactoryProtocol? { get }

  func `init`(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> UnaryCall<Account_Request, Account_Response>

  func add(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> UnaryCall<Account_Request, Account_Response>

  func get(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> UnaryCall<Account_Request, Account_Response>

  func update(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> UnaryCall<Account_Request, Account_Response>

  func delete(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> UnaryCall<Account_Request, Account_Response>

  func list(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> UnaryCall<Account_Request, Account_Response>
}

extension Account_AddServiceClientProtocol {
  internal var serviceName: String {
    return "account.AddService"
  }

  /// Unary call to Init
  ///
  /// - Parameters:
  ///   - request: Request to send to Init.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func `init`(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Account_Request, Account_Response> {
    return self.makeUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.`init`.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeInitInterceptors() ?? []
    )
  }

  /// Unary call to Add
  ///
  /// - Parameters:
  ///   - request: Request to send to Add.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func add(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Account_Request, Account_Response> {
    return self.makeUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.add.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeAddInterceptors() ?? []
    )
  }

  /// Unary call to Get
  ///
  /// - Parameters:
  ///   - request: Request to send to Get.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func get(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Account_Request, Account_Response> {
    return self.makeUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.get.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeGetInterceptors() ?? []
    )
  }

  /// Unary call to Update
  ///
  /// - Parameters:
  ///   - request: Request to send to Update.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func update(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Account_Request, Account_Response> {
    return self.makeUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.update.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeUpdateInterceptors() ?? []
    )
  }

  /// Unary call to Delete
  ///
  /// - Parameters:
  ///   - request: Request to send to Delete.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func delete(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Account_Request, Account_Response> {
    return self.makeUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.delete.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeDeleteInterceptors() ?? []
    )
  }

  /// Unary call to List
  ///
  /// - Parameters:
  ///   - request: Request to send to List.
  ///   - callOptions: Call options.
  /// - Returns: A `UnaryCall` with futures for the metadata, status and response.
  internal func list(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> UnaryCall<Account_Request, Account_Response> {
    return self.makeUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.list.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeListInterceptors() ?? []
    )
  }
}

@available(*, deprecated)
extension Account_AddServiceClient: @unchecked Sendable {}

@available(*, deprecated, renamed: "Account_AddServiceNIOClient")
internal final class Account_AddServiceClient: Account_AddServiceClientProtocol {
  private let lock = Lock()
  private var _defaultCallOptions: CallOptions
  private var _interceptors: Account_AddServiceClientInterceptorFactoryProtocol?
  internal let channel: GRPCChannel
  internal var defaultCallOptions: CallOptions {
    get { self.lock.withLock { return self._defaultCallOptions } }
    set { self.lock.withLockVoid { self._defaultCallOptions = newValue } }
  }
  internal var interceptors: Account_AddServiceClientInterceptorFactoryProtocol? {
    get { self.lock.withLock { return self._interceptors } }
    set { self.lock.withLockVoid { self._interceptors = newValue } }
  }

  /// Creates a client for the account.AddService service.
  ///
  /// - Parameters:
  ///   - channel: `GRPCChannel` to the service host.
  ///   - defaultCallOptions: Options to use for each service call if the user doesn't provide them.
  ///   - interceptors: A factory providing interceptors for each RPC.
  internal init(
    channel: GRPCChannel,
    defaultCallOptions: CallOptions = CallOptions(),
    interceptors: Account_AddServiceClientInterceptorFactoryProtocol? = nil
  ) {
    self.channel = channel
    self._defaultCallOptions = defaultCallOptions
    self._interceptors = interceptors
  }
}

internal struct Account_AddServiceNIOClient: Account_AddServiceClientProtocol {
  internal var channel: GRPCChannel
  internal var defaultCallOptions: CallOptions
  internal var interceptors: Account_AddServiceClientInterceptorFactoryProtocol?

  /// Creates a client for the account.AddService service.
  ///
  /// - Parameters:
  ///   - channel: `GRPCChannel` to the service host.
  ///   - defaultCallOptions: Options to use for each service call if the user doesn't provide them.
  ///   - interceptors: A factory providing interceptors for each RPC.
  internal init(
    channel: GRPCChannel,
    defaultCallOptions: CallOptions = CallOptions(),
    interceptors: Account_AddServiceClientInterceptorFactoryProtocol? = nil
  ) {
    self.channel = channel
    self.defaultCallOptions = defaultCallOptions
    self.interceptors = interceptors
  }
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
internal protocol Account_AddServiceAsyncClientProtocol: GRPCClient {
  static var serviceDescriptor: GRPCServiceDescriptor { get }
  var interceptors: Account_AddServiceClientInterceptorFactoryProtocol? { get }

  func makeInitCall(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response>

  func makeAddCall(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response>

  func makeGetCall(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response>

  func makeUpdateCall(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response>

  func makeDeleteCall(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response>

  func makeListCall(
    _ request: Account_Request,
    callOptions: CallOptions?
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response>
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
extension Account_AddServiceAsyncClientProtocol {
  internal static var serviceDescriptor: GRPCServiceDescriptor {
    return Account_AddServiceClientMetadata.serviceDescriptor
  }

  internal var interceptors: Account_AddServiceClientInterceptorFactoryProtocol? {
    return nil
  }

  internal func makeInitCall(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response> {
    return self.makeAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.`init`.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeInitInterceptors() ?? []
    )
  }

  internal func makeAddCall(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response> {
    return self.makeAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.add.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeAddInterceptors() ?? []
    )
  }

  internal func makeGetCall(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response> {
    return self.makeAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.get.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeGetInterceptors() ?? []
    )
  }

  internal func makeUpdateCall(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response> {
    return self.makeAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.update.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeUpdateInterceptors() ?? []
    )
  }

  internal func makeDeleteCall(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response> {
    return self.makeAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.delete.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeDeleteInterceptors() ?? []
    )
  }

  internal func makeListCall(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) -> GRPCAsyncUnaryCall<Account_Request, Account_Response> {
    return self.makeAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.list.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeListInterceptors() ?? []
    )
  }
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
extension Account_AddServiceAsyncClientProtocol {
  internal func `init`(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) async throws -> Account_Response {
    return try await self.performAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.`init`.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeInitInterceptors() ?? []
    )
  }

  internal func add(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) async throws -> Account_Response {
    return try await self.performAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.add.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeAddInterceptors() ?? []
    )
  }

  internal func get(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) async throws -> Account_Response {
    return try await self.performAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.get.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeGetInterceptors() ?? []
    )
  }

  internal func update(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) async throws -> Account_Response {
    return try await self.performAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.update.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeUpdateInterceptors() ?? []
    )
  }

  internal func delete(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) async throws -> Account_Response {
    return try await self.performAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.delete.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeDeleteInterceptors() ?? []
    )
  }

  internal func list(
    _ request: Account_Request,
    callOptions: CallOptions? = nil
  ) async throws -> Account_Response {
    return try await self.performAsyncUnaryCall(
      path: Account_AddServiceClientMetadata.Methods.list.path,
      request: request,
      callOptions: callOptions ?? self.defaultCallOptions,
      interceptors: self.interceptors?.makeListInterceptors() ?? []
    )
  }
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
internal struct Account_AddServiceAsyncClient: Account_AddServiceAsyncClientProtocol {
  internal var channel: GRPCChannel
  internal var defaultCallOptions: CallOptions
  internal var interceptors: Account_AddServiceClientInterceptorFactoryProtocol?

  internal init(
    channel: GRPCChannel,
    defaultCallOptions: CallOptions = CallOptions(),
    interceptors: Account_AddServiceClientInterceptorFactoryProtocol? = nil
  ) {
    self.channel = channel
    self.defaultCallOptions = defaultCallOptions
    self.interceptors = interceptors
  }
}

internal protocol Account_AddServiceClientInterceptorFactoryProtocol: Sendable {

  /// - Returns: Interceptors to use when invoking '`init`'.
  func makeInitInterceptors() -> [ClientInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when invoking 'add'.
  func makeAddInterceptors() -> [ClientInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when invoking 'get'.
  func makeGetInterceptors() -> [ClientInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when invoking 'update'.
  func makeUpdateInterceptors() -> [ClientInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when invoking 'delete'.
  func makeDeleteInterceptors() -> [ClientInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when invoking 'list'.
  func makeListInterceptors() -> [ClientInterceptor<Account_Request, Account_Response>]
}

internal enum Account_AddServiceClientMetadata {
  internal static let serviceDescriptor = GRPCServiceDescriptor(
    name: "AddService",
    fullName: "account.AddService",
    methods: [
      Account_AddServiceClientMetadata.Methods.`init`,
      Account_AddServiceClientMetadata.Methods.add,
      Account_AddServiceClientMetadata.Methods.get,
      Account_AddServiceClientMetadata.Methods.update,
      Account_AddServiceClientMetadata.Methods.delete,
      Account_AddServiceClientMetadata.Methods.list,
    ]
  )

  internal enum Methods {
    internal static let `init` = GRPCMethodDescriptor(
      name: "Init",
      path: "/account.AddService/Init",
      type: GRPCCallType.unary
    )

    internal static let add = GRPCMethodDescriptor(
      name: "Add",
      path: "/account.AddService/Add",
      type: GRPCCallType.unary
    )

    internal static let get = GRPCMethodDescriptor(
      name: "Get",
      path: "/account.AddService/Get",
      type: GRPCCallType.unary
    )

    internal static let update = GRPCMethodDescriptor(
      name: "Update",
      path: "/account.AddService/Update",
      type: GRPCCallType.unary
    )

    internal static let delete = GRPCMethodDescriptor(
      name: "Delete",
      path: "/account.AddService/Delete",
      type: GRPCCallType.unary
    )

    internal static let list = GRPCMethodDescriptor(
      name: "List",
      path: "/account.AddService/List",
      type: GRPCCallType.unary
    )
  }
}

/// To build a server, implement a class that conforms to this protocol.
internal protocol Account_AddServiceProvider: CallHandlerProvider {
  var interceptors: Account_AddServiceServerInterceptorFactoryProtocol? { get }

  func `init`(request: Account_Request, context: StatusOnlyCallContext) -> EventLoopFuture<Account_Response>

  func add(request: Account_Request, context: StatusOnlyCallContext) -> EventLoopFuture<Account_Response>

  func get(request: Account_Request, context: StatusOnlyCallContext) -> EventLoopFuture<Account_Response>

  func update(request: Account_Request, context: StatusOnlyCallContext) -> EventLoopFuture<Account_Response>

  func delete(request: Account_Request, context: StatusOnlyCallContext) -> EventLoopFuture<Account_Response>

  func list(request: Account_Request, context: StatusOnlyCallContext) -> EventLoopFuture<Account_Response>
}

extension Account_AddServiceProvider {
  internal var serviceName: Substring {
    return Account_AddServiceServerMetadata.serviceDescriptor.fullName[...]
  }

  /// Determines, calls and returns the appropriate request handler, depending on the request's method.
  /// Returns nil for methods not handled by this service.
  internal func handle(
    method name: Substring,
    context: CallHandlerContext
  ) -> GRPCServerHandlerProtocol? {
    switch name {
    case "Init":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeInitInterceptors() ?? [],
        userFunction: self.`init`(request:context:)
      )

    case "Add":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeAddInterceptors() ?? [],
        userFunction: self.add(request:context:)
      )

    case "Get":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeGetInterceptors() ?? [],
        userFunction: self.get(request:context:)
      )

    case "Update":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeUpdateInterceptors() ?? [],
        userFunction: self.update(request:context:)
      )

    case "Delete":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeDeleteInterceptors() ?? [],
        userFunction: self.delete(request:context:)
      )

    case "List":
      return UnaryServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeListInterceptors() ?? [],
        userFunction: self.list(request:context:)
      )

    default:
      return nil
    }
  }
}

/// To implement a server, implement an object which conforms to this protocol.
@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
internal protocol Account_AddServiceAsyncProvider: CallHandlerProvider, Sendable {
  static var serviceDescriptor: GRPCServiceDescriptor { get }
  var interceptors: Account_AddServiceServerInterceptorFactoryProtocol? { get }

  func `init`(
    request: Account_Request,
    context: GRPCAsyncServerCallContext
  ) async throws -> Account_Response

  func add(
    request: Account_Request,
    context: GRPCAsyncServerCallContext
  ) async throws -> Account_Response

  func get(
    request: Account_Request,
    context: GRPCAsyncServerCallContext
  ) async throws -> Account_Response

  func update(
    request: Account_Request,
    context: GRPCAsyncServerCallContext
  ) async throws -> Account_Response

  func delete(
    request: Account_Request,
    context: GRPCAsyncServerCallContext
  ) async throws -> Account_Response

  func list(
    request: Account_Request,
    context: GRPCAsyncServerCallContext
  ) async throws -> Account_Response
}

@available(macOS 10.15, iOS 13, tvOS 13, watchOS 6, *)
extension Account_AddServiceAsyncProvider {
  internal static var serviceDescriptor: GRPCServiceDescriptor {
    return Account_AddServiceServerMetadata.serviceDescriptor
  }

  internal var serviceName: Substring {
    return Account_AddServiceServerMetadata.serviceDescriptor.fullName[...]
  }

  internal var interceptors: Account_AddServiceServerInterceptorFactoryProtocol? {
    return nil
  }

  internal func handle(
    method name: Substring,
    context: CallHandlerContext
  ) -> GRPCServerHandlerProtocol? {
    switch name {
    case "Init":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeInitInterceptors() ?? [],
        wrapping: { try await self.`init`(request: $0, context: $1) }
      )

    case "Add":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeAddInterceptors() ?? [],
        wrapping: { try await self.add(request: $0, context: $1) }
      )

    case "Get":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeGetInterceptors() ?? [],
        wrapping: { try await self.get(request: $0, context: $1) }
      )

    case "Update":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeUpdateInterceptors() ?? [],
        wrapping: { try await self.update(request: $0, context: $1) }
      )

    case "Delete":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeDeleteInterceptors() ?? [],
        wrapping: { try await self.delete(request: $0, context: $1) }
      )

    case "List":
      return GRPCAsyncServerHandler(
        context: context,
        requestDeserializer: ProtobufDeserializer<Account_Request>(),
        responseSerializer: ProtobufSerializer<Account_Response>(),
        interceptors: self.interceptors?.makeListInterceptors() ?? [],
        wrapping: { try await self.list(request: $0, context: $1) }
      )

    default:
      return nil
    }
  }
}

internal protocol Account_AddServiceServerInterceptorFactoryProtocol: Sendable {

  /// - Returns: Interceptors to use when handling '`init`'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeInitInterceptors() -> [ServerInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when handling 'add'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeAddInterceptors() -> [ServerInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when handling 'get'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeGetInterceptors() -> [ServerInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when handling 'update'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeUpdateInterceptors() -> [ServerInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when handling 'delete'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeDeleteInterceptors() -> [ServerInterceptor<Account_Request, Account_Response>]

  /// - Returns: Interceptors to use when handling 'list'.
  ///   Defaults to calling `self.makeInterceptors()`.
  func makeListInterceptors() -> [ServerInterceptor<Account_Request, Account_Response>]
}

internal enum Account_AddServiceServerMetadata {
  internal static let serviceDescriptor = GRPCServiceDescriptor(
    name: "AddService",
    fullName: "account.AddService",
    methods: [
      Account_AddServiceServerMetadata.Methods.`init`,
      Account_AddServiceServerMetadata.Methods.add,
      Account_AddServiceServerMetadata.Methods.get,
      Account_AddServiceServerMetadata.Methods.update,
      Account_AddServiceServerMetadata.Methods.delete,
      Account_AddServiceServerMetadata.Methods.list,
    ]
  )

  internal enum Methods {
    internal static let `init` = GRPCMethodDescriptor(
      name: "Init",
      path: "/account.AddService/Init",
      type: GRPCCallType.unary
    )

    internal static let add = GRPCMethodDescriptor(
      name: "Add",
      path: "/account.AddService/Add",
      type: GRPCCallType.unary
    )

    internal static let get = GRPCMethodDescriptor(
      name: "Get",
      path: "/account.AddService/Get",
      type: GRPCCallType.unary
    )

    internal static let update = GRPCMethodDescriptor(
      name: "Update",
      path: "/account.AddService/Update",
      type: GRPCCallType.unary
    )

    internal static let delete = GRPCMethodDescriptor(
      name: "Delete",
      path: "/account.AddService/Delete",
      type: GRPCCallType.unary
    )

    internal static let list = GRPCMethodDescriptor(
      name: "List",
      path: "/account.AddService/List",
      type: GRPCCallType.unary
    )
  }
}