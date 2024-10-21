// Code generated by protoc-gen-go-vtproto. DO NOT EDIT.
// protoc-gen-go-vtproto version: v0.6.1-0.20240409071808-615f978279ca
// source: google/type/timeofday.proto

package _type

import (
	fmt "fmt"
	protohelpers "github.com/planetscale/vtprotobuf/protohelpers"
	proto "google.golang.org/protobuf/proto"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

func (m *TimeOfDay) CloneVT() *TimeOfDay {
	if m == nil {
		return (*TimeOfDay)(nil)
	}
	r := new(TimeOfDay)
	r.Hours = m.Hours
	r.Minutes = m.Minutes
	r.Seconds = m.Seconds
	r.Nanos = m.Nanos
	if len(m.unknownFields) > 0 {
		r.unknownFields = make([]byte, len(m.unknownFields))
		copy(r.unknownFields, m.unknownFields)
	}
	return r
}

func (m *TimeOfDay) CloneMessageVT() proto.Message {
	return m.CloneVT()
}

func (this *TimeOfDay) EqualVT(that *TimeOfDay) bool {
	if this == that {
		return true
	} else if this == nil || that == nil {
		return false
	}
	if this.Hours != that.Hours {
		return false
	}
	if this.Minutes != that.Minutes {
		return false
	}
	if this.Seconds != that.Seconds {
		return false
	}
	if this.Nanos != that.Nanos {
		return false
	}
	return string(this.unknownFields) == string(that.unknownFields)
}

func (this *TimeOfDay) EqualMessageVT(thatMsg proto.Message) bool {
	that, ok := thatMsg.(*TimeOfDay)
	if !ok {
		return false
	}
	return this.EqualVT(that)
}
func (m *TimeOfDay) MarshalVT() (dAtA []byte, err error) {
	if m == nil {
		return nil, nil
	}
	size := m.SizeVT()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBufferVT(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TimeOfDay) MarshalToVT(dAtA []byte) (int, error) {
	size := m.SizeVT()
	return m.MarshalToSizedBufferVT(dAtA[:size])
}

func (m *TimeOfDay) MarshalToSizedBufferVT(dAtA []byte) (int, error) {
	if m == nil {
		return 0, nil
	}
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.unknownFields != nil {
		i -= len(m.unknownFields)
		copy(dAtA[i:], m.unknownFields)
	}
	if m.Nanos != 0 {
		i = protohelpers.EncodeVarint(dAtA, i, uint64(m.Nanos))
		i--
		dAtA[i] = 0x20
	}
	if m.Seconds != 0 {
		i = protohelpers.EncodeVarint(dAtA, i, uint64(m.Seconds))
		i--
		dAtA[i] = 0x18
	}
	if m.Minutes != 0 {
		i = protohelpers.EncodeVarint(dAtA, i, uint64(m.Minutes))
		i--
		dAtA[i] = 0x10
	}
	if m.Hours != 0 {
		i = protohelpers.EncodeVarint(dAtA, i, uint64(m.Hours))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TimeOfDay) SizeVT() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Hours != 0 {
		n += 1 + protohelpers.SizeOfVarint(uint64(m.Hours))
	}
	if m.Minutes != 0 {
		n += 1 + protohelpers.SizeOfVarint(uint64(m.Minutes))
	}
	if m.Seconds != 0 {
		n += 1 + protohelpers.SizeOfVarint(uint64(m.Seconds))
	}
	if m.Nanos != 0 {
		n += 1 + protohelpers.SizeOfVarint(uint64(m.Nanos))
	}
	n += len(m.unknownFields)
	return n
}

func (m *TimeOfDay) UnmarshalVT(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return protohelpers.ErrIntOverflow
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TimeOfDay: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TimeOfDay: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hours", wireType)
			}
			m.Hours = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protohelpers.ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Hours |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Minutes", wireType)
			}
			m.Minutes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protohelpers.ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Minutes |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Seconds", wireType)
			}
			m.Seconds = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protohelpers.ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Seconds |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nanos", wireType)
			}
			m.Nanos = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protohelpers.ErrIntOverflow
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nanos |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := protohelpers.Skip(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return protohelpers.ErrInvalidLength
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.unknownFields = append(m.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
