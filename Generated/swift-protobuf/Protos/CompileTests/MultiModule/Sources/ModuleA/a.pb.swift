// DO NOT EDIT.
// swift-format-ignore-file
// swiftlint:disable all
//
// Generated by the Swift generator plugin for the protocol buffer compiler.
// Source: swift-protobuf/Protos/CompileTests/MultiModule/Sources/ModuleA/a.proto
//
// For information on using the generated types, please see the documentation:
//   https://github.com/apple/swift-protobuf/

import SwiftProtobuf

// If the compiler emits an error on this type, it is because this file
// was generated by a version of the `protoc` Swift plug-in that is
// incompatible with the version of SwiftProtobuf to which you are linking.
// Please ensure that you are building against the same version of the API
// that was used to generate this file.
fileprivate struct _GeneratedWithProtocGenSwiftVersion: SwiftProtobuf.ProtobufAPIVersionCheck {
  struct _2: SwiftProtobuf.ProtobufAPIVersion_2 {}
  typealias Version = _2
}

enum E: Int, SwiftProtobuf.Enum, Swift.CaseIterable {
  case unset = 0
  case a = 1
  case b = 2

  init() {
    self = .unset
  }

}

struct A: SwiftProtobuf.ExtensibleMessage, Sendable {
  // SwiftProtobuf.Message conformance is added in an extension below. See the
  // `Message` and `Message+*Additions` files in the SwiftProtobuf library for
  // methods supported on all messages.

  var e: E {
    get {return _e ?? .unset}
    set {_e = newValue}
  }
  /// Returns true if `e` has been explicitly set.
  var hasE: Bool {return self._e != nil}
  /// Clears the value of `e`. Subsequent reads from it will return its default value.
  mutating func clearE() {self._e = nil}

  var unknownFields = SwiftProtobuf.UnknownStorage()

  init() {}

  var _protobuf_extensionFieldValues = SwiftProtobuf.ExtensionFieldValueSet()
  fileprivate var _e: E? = nil
}

// MARK: - Extension support defined in a.proto.

// MARK: - Extension Properties

// Swift Extensions on the extended Messages to add easy access to the declared
// extension fields. The names are based on the extension field name from the proto
// declaration. To avoid naming collisions, the names are prefixed with the name of
// the scope where the extend directive occurs.

extension A {

  var extStr: String {
    get {return getExtensionValue(ext: Extensions_ext_str) ?? String()}
    set {setExtensionValue(ext: Extensions_ext_str, value: newValue)}
  }
  /// Returns true if extension `Extensions_ext_str`
  /// has been explicitly set.
  var hasExtStr: Bool {
    return hasExtensionValue(ext: Extensions_ext_str)
  }
  /// Clears the value of extension `Extensions_ext_str`.
  /// Subsequent reads from it will return its default value.
  mutating func clearExtStr() {
    clearExtensionValue(ext: Extensions_ext_str)
  }

}

// MARK: - File's ExtensionMap: A_Extensions

/// A `SwiftProtobuf.SimpleExtensionMap` that includes all of the extensions defined by
/// this .proto file. It can be used any place an `SwiftProtobuf.ExtensionMap` is needed
/// in parsing, or it can be combined with other `SwiftProtobuf.SimpleExtensionMap`s to create
/// a larger `SwiftProtobuf.SimpleExtensionMap`.
let A_Extensions: SwiftProtobuf.SimpleExtensionMap = [
  Extensions_ext_str
]

// Extension Objects - The only reason these might be needed is when manually
// constructing a `SimpleExtensionMap`, otherwise, use the above _Extension Properties_
// accessors for the extension fields on the messages directly.

let Extensions_ext_str = SwiftProtobuf.MessageExtension<SwiftProtobuf.OptionalExtensionField<SwiftProtobuf.ProtobufString>, A>(
  _protobuf_fieldNumber: 100,
  fieldName: "ext_str"
)

// MARK: - Code below here is support for the SwiftProtobuf runtime.

extension E: SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    0: .same(proto: "E_UNSET"),
    1: .same(proto: "E_A"),
    2: .same(proto: "E_B"),
  ]
}

extension A: SwiftProtobuf.Message, SwiftProtobuf._MessageImplementationBase, SwiftProtobuf._ProtoNameProviding {
  static let protoMessageName: String = "A"
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    1: .same(proto: "e"),
  ]

  public var isInitialized: Bool {
    if !_protobuf_extensionFieldValues.isInitialized {return false}
    return true
  }

  mutating func decodeMessage<D: SwiftProtobuf.Decoder>(decoder: inout D) throws {
    while let fieldNumber = try decoder.nextFieldNumber() {
      // The use of inline closures is to circumvent an issue where the compiler
      // allocates stack space for every case branch when no optimizations are
      // enabled. https://github.com/apple/swift-protobuf/issues/1034
      switch fieldNumber {
      case 1: try { try decoder.decodeSingularEnumField(value: &self._e) }()
      case 100..<1001:
        try { try decoder.decodeExtensionField(values: &_protobuf_extensionFieldValues, messageType: A.self, fieldNumber: fieldNumber) }()
      default: break
      }
    }
  }

  func traverse<V: SwiftProtobuf.Visitor>(visitor: inout V) throws {
    // The use of inline closures is to circumvent an issue where the compiler
    // allocates stack space for every if/case branch local when no optimizations
    // are enabled. https://github.com/apple/swift-protobuf/issues/1034 and
    // https://github.com/apple/swift-protobuf/issues/1182
    try { if let v = self._e {
      try visitor.visitSingularEnumField(value: v, fieldNumber: 1)
    } }()
    try visitor.visitExtensionFields(fields: _protobuf_extensionFieldValues, start: 100, end: 1001)
    try unknownFields.traverse(visitor: &visitor)
  }

  static func ==(lhs: A, rhs: A) -> Bool {
    if lhs._e != rhs._e {return false}
    if lhs.unknownFields != rhs.unknownFields {return false}
    if lhs._protobuf_extensionFieldValues != rhs._protobuf_extensionFieldValues {return false}
    return true
  }
}
