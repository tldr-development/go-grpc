// DO NOT EDIT.
// swift-format-ignore-file
// swiftlint:disable all
//
// Generated by the Swift generator plugin for the protocol buffer compiler.
// Source: swift-protobuf/Protos/upstream/google/protobuf/unittest_proto3_bad_macros.proto
//
// For information on using the generated types, please see the documentation:
//   https://github.com/apple/swift-protobuf/

// Protocol Buffers - Google's data interchange format
// Copyright 2023 Google Inc.  All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

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

/// This generates `GID_MAX`, which is a macro in some circumstances.
enum ProtobufUnittest_GID: SwiftProtobuf.Enum, Swift.CaseIterable {
  typealias RawValue = Int
  case unused // = 0
  case UNRECOGNIZED(Int)

  init() {
    self = .unused
  }

  init?(rawValue: Int) {
    switch rawValue {
    case 0: self = .unused
    default: self = .UNRECOGNIZED(rawValue)
    }
  }

  var rawValue: Int {
    switch self {
    case .unused: return 0
    case .UNRECOGNIZED(let i): return i
    }
  }

  // The compiler won't synthesize support with the UNRECOGNIZED case.
  static let allCases: [ProtobufUnittest_GID] = [
    .unused,
  ]

}

/// This generates `UID_MAX`, which is a mcro in some circumstances.
enum ProtobufUnittest_UID: SwiftProtobuf.Enum, Swift.CaseIterable {
  typealias RawValue = Int
  case unused // = 0
  case UNRECOGNIZED(Int)

  init() {
    self = .unused
  }

  init?(rawValue: Int) {
    switch rawValue {
    case 0: self = .unused
    default: self = .UNRECOGNIZED(rawValue)
    }
  }

  var rawValue: Int {
    switch self {
    case .unused: return 0
    case .UNRECOGNIZED(let i): return i
    }
  }

  // The compiler won't synthesize support with the UNRECOGNIZED case.
  static let allCases: [ProtobufUnittest_UID] = [
    .unused,
  ]

}

/// Just a container for bad macro names. Some of these do not follow the normal
/// naming conventions, this is intentional, we just want to trigger a build
/// failure if the macro is left defined.
enum ProtobufUnittest_BadNames: SwiftProtobuf.Enum, Swift.CaseIterable {
  typealias RawValue = Int

  /// autoheader defines this in some circumstances.
  case package // = 0

  /// The comment says "a few common headers define this".
  case packed // = 1

  /// Defined in many Linux system headers.
  case linux // = 2

  /// This is often a macro in `<math.h>`.
  case domain // = 3

  /// These are defined in both Windows and macOS headers.
  case `true` // = 4
  case `false` // = 5

  /// Sometimes defined in Windows system headers.
  case createNew // = 6
  case delete // = 7
  case doubleClick // = 8
  case error // = 9
  case errorBusy // = 10
  case errorInstallFailed // = 11
  case errorNotFound // = 12
  case getClassName // = 13
  case getCurrentTime // = 14
  case getMessage // = 15
  case getObject // = 16
  case ignore // = 17
  case `in` // = 18
  case inputKeyboard // = 19
  case noError // = 20
  case out // = 21
  case `optional` // = 22
  case near // = 23
  case noData // = 24
  case reasonUnknown // = 25
  case serviceDisabled // = 26
  case severityError // = 27
  case statusPending // = 28
  case strict // = 29

  /// Sometimed defined in macOS system headers.
  case typeBool // = 30

  /// Defined in macOS, Windows, and Linux headers.
  case debug // = 31
  case UNRECOGNIZED(Int)

  init() {
    self = .package
  }

  init?(rawValue: Int) {
    switch rawValue {
    case 0: self = .package
    case 1: self = .packed
    case 2: self = .linux
    case 3: self = .domain
    case 4: self = .true
    case 5: self = .false
    case 6: self = .createNew
    case 7: self = .delete
    case 8: self = .doubleClick
    case 9: self = .error
    case 10: self = .errorBusy
    case 11: self = .errorInstallFailed
    case 12: self = .errorNotFound
    case 13: self = .getClassName
    case 14: self = .getCurrentTime
    case 15: self = .getMessage
    case 16: self = .getObject
    case 17: self = .ignore
    case 18: self = .in
    case 19: self = .inputKeyboard
    case 20: self = .noError
    case 21: self = .out
    case 22: self = .optional
    case 23: self = .near
    case 24: self = .noData
    case 25: self = .reasonUnknown
    case 26: self = .serviceDisabled
    case 27: self = .severityError
    case 28: self = .statusPending
    case 29: self = .strict
    case 30: self = .typeBool
    case 31: self = .debug
    default: self = .UNRECOGNIZED(rawValue)
    }
  }

  var rawValue: Int {
    switch self {
    case .package: return 0
    case .packed: return 1
    case .linux: return 2
    case .domain: return 3
    case .true: return 4
    case .false: return 5
    case .createNew: return 6
    case .delete: return 7
    case .doubleClick: return 8
    case .error: return 9
    case .errorBusy: return 10
    case .errorInstallFailed: return 11
    case .errorNotFound: return 12
    case .getClassName: return 13
    case .getCurrentTime: return 14
    case .getMessage: return 15
    case .getObject: return 16
    case .ignore: return 17
    case .in: return 18
    case .inputKeyboard: return 19
    case .noError: return 20
    case .out: return 21
    case .optional: return 22
    case .near: return 23
    case .noData: return 24
    case .reasonUnknown: return 25
    case .serviceDisabled: return 26
    case .severityError: return 27
    case .statusPending: return 28
    case .strict: return 29
    case .typeBool: return 30
    case .debug: return 31
    case .UNRECOGNIZED(let i): return i
    }
  }

  // The compiler won't synthesize support with the UNRECOGNIZED case.
  static let allCases: [ProtobufUnittest_BadNames] = [
    .package,
    .packed,
    .linux,
    .domain,
    .true,
    .false,
    .createNew,
    .delete,
    .doubleClick,
    .error,
    .errorBusy,
    .errorInstallFailed,
    .errorNotFound,
    .getClassName,
    .getCurrentTime,
    .getMessage,
    .getObject,
    .ignore,
    .in,
    .inputKeyboard,
    .noError,
    .out,
    .optional,
    .near,
    .noData,
    .reasonUnknown,
    .serviceDisabled,
    .severityError,
    .statusPending,
    .strict,
    .typeBool,
    .debug,
  ]

}

// MARK: - Code below here is support for the SwiftProtobuf runtime.

extension ProtobufUnittest_GID: SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    0: .same(proto: "GID_UNUSED"),
  ]
}

extension ProtobufUnittest_UID: SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    0: .same(proto: "UID_UNUSED"),
  ]
}

extension ProtobufUnittest_BadNames: SwiftProtobuf._ProtoNameProviding {
  static let _protobuf_nameMap: SwiftProtobuf._NameMap = [
    0: .same(proto: "PACKAGE"),
    1: .same(proto: "PACKED"),
    2: .same(proto: "linux"),
    3: .same(proto: "DOMAIN"),
    4: .same(proto: "TRUE"),
    5: .same(proto: "FALSE"),
    6: .same(proto: "CREATE_NEW"),
    7: .same(proto: "DELETE"),
    8: .same(proto: "DOUBLE_CLICK"),
    9: .same(proto: "ERROR"),
    10: .same(proto: "ERROR_BUSY"),
    11: .same(proto: "ERROR_INSTALL_FAILED"),
    12: .same(proto: "ERROR_NOT_FOUND"),
    13: .same(proto: "GetClassName"),
    14: .same(proto: "GetCurrentTime"),
    15: .same(proto: "GetMessage"),
    16: .same(proto: "GetObject"),
    17: .same(proto: "IGNORE"),
    18: .same(proto: "IN"),
    19: .same(proto: "INPUT_KEYBOARD"),
    20: .same(proto: "NO_ERROR"),
    21: .same(proto: "OUT"),
    22: .same(proto: "OPTIONAL"),
    23: .same(proto: "NEAR"),
    24: .same(proto: "NO_DATA"),
    25: .same(proto: "REASON_UNKNOWN"),
    26: .same(proto: "SERVICE_DISABLED"),
    27: .same(proto: "SEVERITY_ERROR"),
    28: .same(proto: "STATUS_PENDING"),
    29: .same(proto: "STRICT"),
    30: .same(proto: "TYPE_BOOL"),
    31: .same(proto: "DEBUG"),
  ]
}
