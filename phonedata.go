package phonedata

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

const (
	CMCC    byte = iota + 0x01 //中国移动
	CUCC                       //中国联通
	CTCC                       //中国电信
	CTCCVNO                    //电信虚拟运营商
	CUCCVNO                    //联通虚拟运营商
	CMCCVNO                    //移动虚拟运营商
)

const (
	IntLen         uint32 = 4
	PhoneEntrySize uint32 = 9
	UnkownOperator        = "未知电信运营商"
)

type PhoneRecord struct {
	PhoneNum string
	Province string
	City     string
	ZipCode  string
	AreaZone string
	CardType string
}

var (
	// 1 移动
	// 2 联通
	// 3 电信
	// 4 电信虚拟运营商
	// 5 联通虚拟运营商
	// 6 移动虚拟运营商
	CardTypeMap = map[byte]string{
		CMCC:    "中国移动",
		CUCC:    "中国联通",
		CTCC:    "中国电信",
		CTCCVNO: "中国电信虚拟运营商",
		CUCCVNO: "中国联通虚拟运营商",
		CMCCVNO: "中国移动虚拟运营商",
	}
)

// phone.dat文件格式
//         | 4 bytes |                     <- phone.dat 版本号
//         ------------
//         | 4 bytes |                     <-  第一个索引的偏移
//         -----------------------
//         |  offset - 8            |      <-  记录区
//         -----------------------
//         |  index                 |      <-  索引区
//         -----------------------
// 小端存储
// 头部 头部为8个字节，版本号为4个字节，第一个索引的偏移为4个字节(<4si)。
// 记录区 中每条记录的格式为"<省份>|<城市>|<邮编>|<长途区号>\0"。 每条记录以'\0'结束。
// 索引区 中每条记录的格式为"<手机号前七位><记录区的偏移><卡类型>"，每个索引的长度为9个字节(<iiB)。
// 解析步骤:

// 解析头部8个字节，得到索引区的第一条索引的偏移。
// 在索引区用二分查找得出手机号在记录区的记录偏移。
// 在记录区从上一步得到的记录偏移处取数据，直到遇到'\0'。
// 我定义的卡类型为:

// Parser is phone parser
type Parser struct {
	bin         []byte
	firstoffset uint32
}

// NewParser create a phone parser
func NewParser(r io.Reader) (*Parser, error) {
	bin, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	p := &Parser{
		bin: bin,
	}
	p.firstoffset = p.firstRecordOffset()

	return p, nil
}

// Version return phone.dat version number
func (p *Parser) Version() string {
	return string(p.bin[:4])
}

func (p *Parser) firstRecordOffset() uint32 {
	return binary.LittleEndian.Uint32(p.bin[4:8])
}

func (p *Parser) totalRecord() uint32 {
	return (uint32(len(p.bin)) - p.firstRecordOffset()) / PhoneEntrySize
}

func (p *Parser) Find(phone string) (*PhoneRecord, error) {
	if len(phone) < 7 || len(phone) > 11 {
		return nil, errors.New("illegal phone length")
	}
	pi, err := strconv.Atoi(phone[:7])
	if err != nil {
		return nil, err
	}
	target := uint32(pi)

	var left uint32
	right := p.totalRecord()
	for left < right {
		mid := left + (right-left)/2
		offset := p.firstoffset + mid*PhoneEntrySize
		// 索引区 中每条记录的格式为"<手机号前七位><记录区的偏移><卡类型>"，每个索引的长度为9个字节(<iiB)。
		curPhone := binary.LittleEndian.Uint32(p.bin[offset : offset+IntLen])
		if curPhone > target {
			right = mid
		} else if curPhone < target {
			left = mid + 1
		} else {
			recordOffset := int32(binary.LittleEndian.Uint32(p.bin[offset+IntLen : offset+IntLen*2]))
			cardType := p.bin[offset+IntLen*2]
			cbyte := p.bin[recordOffset:]
			// 记录区中 每条记录以'\0'结束。
			endOffset := int32(bytes.Index(cbyte, []byte("\000")))
			data := bytes.Split(cbyte[:endOffset], []byte("|"))
			cardStr, ok := CardTypeMap[cardType]
			if !ok {
				cardStr = UnkownOperator
			}
			// 记录区 中每条记录的格式为"<省份>|<城市>|<邮编>|<长途区号>\0"。 每条记录以'\0'结束。
			pr := &PhoneRecord{
				PhoneNum: phone,
				Province: string(data[0]),
				City:     string(data[1]),
				ZipCode:  string(data[2]),
				AreaZone: string(data[3]),
				CardType: cardStr,
			}
			return pr, nil
		}
	}
	return nil, errors.New("phone's data not found")
}

func (pr *PhoneRecord) String() string {
	return fmt.Sprintf("手机号: %s\n区号: %s\n运营商: %s\n城市: %s\n邮编: %s\n省份: %s\n", pr.PhoneNum, pr.AreaZone, pr.CardType, pr.City, pr.ZipCode, pr.Province)
}
