// Copyright (c) 2019-2020 The Zcash developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or https://www.opensource.org/licenses/mit-license.php .

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: compact_formats.proto

package walletrpc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CompactBlock is a packaging of ONLY the data from a block that's needed to:
//   1. Detect a payment to your shielded Sapling address
//   2. Detect a spend of your shielded Sapling notes
//   3. Update your witnesses to generate new Sapling spend proofs.
type CompactBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProtoVersion uint32       `protobuf:"varint,1,opt,name=protoVersion,proto3" json:"protoVersion,omitempty"` // the version of this wire format, for storage
	Height       uint64       `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`             // the height of this block
	Hash         []byte       `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`                  // the ID (hash) of this block, same as in block explorers
	PrevHash     []byte       `protobuf:"bytes,4,opt,name=prevHash,proto3" json:"prevHash,omitempty"`          // the ID (hash) of this block's predecessor
	Time         uint32       `protobuf:"varint,5,opt,name=time,proto3" json:"time,omitempty"`                 // Unix epoch time when the block was mined
	Header       []byte       `protobuf:"bytes,6,opt,name=header,proto3" json:"header,omitempty"`              // (hash, prevHash, and time) OR (full header)
	Vtx          []*CompactTx `protobuf:"bytes,7,rep,name=vtx,proto3" json:"vtx,omitempty"`                    // zero or more compact transactions from this block
}

func (x *CompactBlock) Reset() {
	*x = CompactBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_compact_formats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompactBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompactBlock) ProtoMessage() {}

func (x *CompactBlock) ProtoReflect() protoreflect.Message {
	mi := &file_compact_formats_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompactBlock.ProtoReflect.Descriptor instead.
func (*CompactBlock) Descriptor() ([]byte, []int) {
	return file_compact_formats_proto_rawDescGZIP(), []int{0}
}

func (x *CompactBlock) GetProtoVersion() uint32 {
	if x != nil {
		return x.ProtoVersion
	}
	return 0
}

func (x *CompactBlock) GetHeight() uint64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *CompactBlock) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *CompactBlock) GetPrevHash() []byte {
	if x != nil {
		return x.PrevHash
	}
	return nil
}

func (x *CompactBlock) GetTime() uint32 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *CompactBlock) GetHeader() []byte {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *CompactBlock) GetVtx() []*CompactTx {
	if x != nil {
		return x.Vtx
	}
	return nil
}

// CompactTx contains the minimum information for a wallet to know if this transaction
// is relevant to it (either pays to it or spends from it) via shielded elements
// only. This message will not encode a transparent-to-transparent transaction.
type CompactTx struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"` // the index within the full block
	Hash  []byte `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`    // the ID (hash) of this transaction, same as in block explorers
	// The transaction fee: present if server can provide. In the case of a
	// stateless server and a transaction with transparent inputs, this will be
	// unset because the calculation requires reference to prior transactions.
	// in a pure-Sapling context, the fee will be calculable as:
	//    valueBalance + (sum(vPubNew) - sum(vPubOld) - sum(tOut))
	Fee     uint32                  `protobuf:"varint,3,opt,name=fee,proto3" json:"fee,omitempty"`
	Spends  []*CompactSaplingSpend  `protobuf:"bytes,4,rep,name=spends,proto3" json:"spends,omitempty"`   // inputs
	Outputs []*CompactSaplingOutput `protobuf:"bytes,5,rep,name=outputs,proto3" json:"outputs,omitempty"` // outputs
	Actions []*CompactOrchardAction `protobuf:"bytes,6,rep,name=actions,proto3" json:"actions,omitempty"`
}

func (x *CompactTx) Reset() {
	*x = CompactTx{}
	if protoimpl.UnsafeEnabled {
		mi := &file_compact_formats_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompactTx) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompactTx) ProtoMessage() {}

func (x *CompactTx) ProtoReflect() protoreflect.Message {
	mi := &file_compact_formats_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompactTx.ProtoReflect.Descriptor instead.
func (*CompactTx) Descriptor() ([]byte, []int) {
	return file_compact_formats_proto_rawDescGZIP(), []int{1}
}

func (x *CompactTx) GetIndex() uint64 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *CompactTx) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *CompactTx) GetFee() uint32 {
	if x != nil {
		return x.Fee
	}
	return 0
}

func (x *CompactTx) GetSpends() []*CompactSaplingSpend {
	if x != nil {
		return x.Spends
	}
	return nil
}

func (x *CompactTx) GetOutputs() []*CompactSaplingOutput {
	if x != nil {
		return x.Outputs
	}
	return nil
}

func (x *CompactTx) GetActions() []*CompactOrchardAction {
	if x != nil {
		return x.Actions
	}
	return nil
}

// CompactSaplingSpend is a Sapling Spend Description as described in 7.3 of the Zcash
// protocol specification.
type CompactSaplingSpend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nf []byte `protobuf:"bytes,1,opt,name=nf,proto3" json:"nf,omitempty"` // nullifier (see the Zcash protocol specification)
}

func (x *CompactSaplingSpend) Reset() {
	*x = CompactSaplingSpend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_compact_formats_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompactSaplingSpend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompactSaplingSpend) ProtoMessage() {}

func (x *CompactSaplingSpend) ProtoReflect() protoreflect.Message {
	mi := &file_compact_formats_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompactSaplingSpend.ProtoReflect.Descriptor instead.
func (*CompactSaplingSpend) Descriptor() ([]byte, []int) {
	return file_compact_formats_proto_rawDescGZIP(), []int{2}
}

func (x *CompactSaplingSpend) GetNf() []byte {
	if x != nil {
		return x.Nf
	}
	return nil
}

// output is a Sapling Output Description as described in section 7.4 of the
// Zcash protocol spec. Total size is 948.
type CompactSaplingOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmu        []byte `protobuf:"bytes,1,opt,name=cmu,proto3" json:"cmu,omitempty"`               // note commitment u-coordinate
	Epk        []byte `protobuf:"bytes,2,opt,name=epk,proto3" json:"epk,omitempty"`               // ephemeral public key
	Ciphertext []byte `protobuf:"bytes,3,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"` // first 52 bytes of ciphertext
}

func (x *CompactSaplingOutput) Reset() {
	*x = CompactSaplingOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_compact_formats_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompactSaplingOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompactSaplingOutput) ProtoMessage() {}

func (x *CompactSaplingOutput) ProtoReflect() protoreflect.Message {
	mi := &file_compact_formats_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompactSaplingOutput.ProtoReflect.Descriptor instead.
func (*CompactSaplingOutput) Descriptor() ([]byte, []int) {
	return file_compact_formats_proto_rawDescGZIP(), []int{3}
}

func (x *CompactSaplingOutput) GetCmu() []byte {
	if x != nil {
		return x.Cmu
	}
	return nil
}

func (x *CompactSaplingOutput) GetEpk() []byte {
	if x != nil {
		return x.Epk
	}
	return nil
}

func (x *CompactSaplingOutput) GetCiphertext() []byte {
	if x != nil {
		return x.Ciphertext
	}
	return nil
}

// https://github.com/zcash/zips/blob/main/zip-0225.rst#orchard-action-description-orchardaction
// (but not all fields are needed)
type CompactOrchardAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nullifier     []byte `protobuf:"bytes,1,opt,name=nullifier,proto3" json:"nullifier,omitempty"`         // [32] The nullifier of the input note
	Cmx           []byte `protobuf:"bytes,2,opt,name=cmx,proto3" json:"cmx,omitempty"`                     // [32] The x-coordinate of the note commitment for the output note
	EphemeralKey  []byte `protobuf:"bytes,3,opt,name=ephemeralKey,proto3" json:"ephemeralKey,omitempty"`   // [32] An encoding of an ephemeral Pallas public key
	EncCiphertext []byte `protobuf:"bytes,4,opt,name=encCiphertext,proto3" json:"encCiphertext,omitempty"` // [580] The encrypted contents of the note plaintext
}

func (x *CompactOrchardAction) Reset() {
	*x = CompactOrchardAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_compact_formats_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompactOrchardAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompactOrchardAction) ProtoMessage() {}

func (x *CompactOrchardAction) ProtoReflect() protoreflect.Message {
	mi := &file_compact_formats_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompactOrchardAction.ProtoReflect.Descriptor instead.
func (*CompactOrchardAction) Descriptor() ([]byte, []int) {
	return file_compact_formats_proto_rawDescGZIP(), []int{4}
}

func (x *CompactOrchardAction) GetNullifier() []byte {
	if x != nil {
		return x.Nullifier
	}
	return nil
}

func (x *CompactOrchardAction) GetCmx() []byte {
	if x != nil {
		return x.Cmx
	}
	return nil
}

func (x *CompactOrchardAction) GetEphemeralKey() []byte {
	if x != nil {
		return x.EphemeralKey
	}
	return nil
}

func (x *CompactOrchardAction) GetEncCiphertext() []byte {
	if x != nil {
		return x.EncCiphertext
	}
	return nil
}

var File_compact_formats_proto protoreflect.FileDescriptor

var file_compact_formats_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x5f, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x63, 0x61, 0x73, 0x68, 0x2e, 0x7a, 0x2e,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x22, 0xda,
	0x01, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12,
	0x22, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x56, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x72, 0x65, 0x76, 0x48, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x08, 0x70, 0x72, 0x65, 0x76, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x03, 0x76, 0x74, 0x78, 0x18, 0x07,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x61, 0x73, 0x68, 0x2e, 0x7a, 0x2e, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f, 0x6d,
	0x70, 0x61, 0x63, 0x74, 0x54, 0x78, 0x52, 0x03, 0x76, 0x74, 0x78, 0x22, 0x99, 0x02, 0x0a, 0x09,
	0x43, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x54, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x65, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x03, 0x66, 0x65, 0x65, 0x12, 0x42, 0x0a, 0x06, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x18,
	0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x63, 0x61, 0x73, 0x68, 0x2e, 0x7a, 0x2e, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f,
	0x6d, 0x70, 0x61, 0x63, 0x74, 0x53, 0x61, 0x70, 0x6c, 0x69, 0x6e, 0x67, 0x53, 0x70, 0x65, 0x6e,
	0x64, 0x52, 0x06, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x73, 0x12, 0x45, 0x0a, 0x07, 0x6f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x61, 0x73,
	0x68, 0x2e, 0x7a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x53, 0x61, 0x70, 0x6c, 0x69, 0x6e,
	0x67, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x07, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x73,
	0x12, 0x45, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x2b, 0x2e, 0x63, 0x61, 0x73, 0x68, 0x2e, 0x7a, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x2e, 0x73, 0x64, 0x6b, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x63,
	0x74, 0x4f, 0x72, 0x63, 0x68, 0x61, 0x72, 0x64, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x25, 0x0a, 0x13, 0x43, 0x6f, 0x6d, 0x70, 0x61,
	0x63, 0x74, 0x53, 0x61, 0x70, 0x6c, 0x69, 0x6e, 0x67, 0x53, 0x70, 0x65, 0x6e, 0x64, 0x12, 0x0e,
	0x0a, 0x02, 0x6e, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x6e, 0x66, 0x22, 0x5a,
	0x0a, 0x14, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x53, 0x61, 0x70, 0x6c, 0x69, 0x6e, 0x67,
	0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6d, 0x75, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x03, 0x63, 0x6d, 0x75, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x70, 0x6b, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x65, 0x70, 0x6b, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x69,
	0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a,
	0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74, 0x22, 0x90, 0x01, 0x0a, 0x14, 0x43,
	0x6f, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x4f, 0x72, 0x63, 0x68, 0x61, 0x72, 0x64, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x75, 0x6c, 0x6c, 0x69, 0x66, 0x69, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x6e, 0x75, 0x6c, 0x6c, 0x69, 0x66, 0x69, 0x65,
	0x72, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x6d, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x63, 0x6d, 0x78, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c,
	0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x65, 0x70, 0x68, 0x65, 0x6d,
	0x65, 0x72, 0x61, 0x6c, 0x4b, 0x65, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x65, 0x6e, 0x63, 0x43, 0x69,
	0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d,
	0x65, 0x6e, 0x63, 0x43, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74, 0x42, 0x1b, 0x5a,
	0x16, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x64, 0x2f, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x72, 0x70, 0x63, 0xba, 0x02, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_compact_formats_proto_rawDescOnce sync.Once
	file_compact_formats_proto_rawDescData = file_compact_formats_proto_rawDesc
)

func file_compact_formats_proto_rawDescGZIP() []byte {
	file_compact_formats_proto_rawDescOnce.Do(func() {
		file_compact_formats_proto_rawDescData = protoimpl.X.CompressGZIP(file_compact_formats_proto_rawDescData)
	})
	return file_compact_formats_proto_rawDescData
}

var file_compact_formats_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_compact_formats_proto_goTypes = []interface{}{
	(*CompactBlock)(nil),         // 0: cash.z.wallet.sdk.rpc.CompactBlock
	(*CompactTx)(nil),            // 1: cash.z.wallet.sdk.rpc.CompactTx
	(*CompactSaplingSpend)(nil),  // 2: cash.z.wallet.sdk.rpc.CompactSaplingSpend
	(*CompactSaplingOutput)(nil), // 3: cash.z.wallet.sdk.rpc.CompactSaplingOutput
	(*CompactOrchardAction)(nil), // 4: cash.z.wallet.sdk.rpc.CompactOrchardAction
}
var file_compact_formats_proto_depIdxs = []int32{
	1, // 0: cash.z.wallet.sdk.rpc.CompactBlock.vtx:type_name -> cash.z.wallet.sdk.rpc.CompactTx
	2, // 1: cash.z.wallet.sdk.rpc.CompactTx.spends:type_name -> cash.z.wallet.sdk.rpc.CompactSaplingSpend
	3, // 2: cash.z.wallet.sdk.rpc.CompactTx.outputs:type_name -> cash.z.wallet.sdk.rpc.CompactSaplingOutput
	4, // 3: cash.z.wallet.sdk.rpc.CompactTx.actions:type_name -> cash.z.wallet.sdk.rpc.CompactOrchardAction
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_compact_formats_proto_init() }
func file_compact_formats_proto_init() {
	if File_compact_formats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_compact_formats_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompactBlock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_compact_formats_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompactTx); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_compact_formats_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompactSaplingSpend); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_compact_formats_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompactSaplingOutput); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_compact_formats_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompactOrchardAction); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_compact_formats_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_compact_formats_proto_goTypes,
		DependencyIndexes: file_compact_formats_proto_depIdxs,
		MessageInfos:      file_compact_formats_proto_msgTypes,
	}.Build()
	File_compact_formats_proto = out.File
	file_compact_formats_proto_rawDesc = nil
	file_compact_formats_proto_goTypes = nil
	file_compact_formats_proto_depIdxs = nil
}
