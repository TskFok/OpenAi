package tool

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"
)

func UUID() string {
	uuidWithHyphen := uuid.New()

	return strings.Replace(uuidWithHyphen.String(), "-", "", -1)
}

func Password(password string, salt string) string {
	builder := strings.Builder{}
	builder.WriteString(password)
	builder.WriteString(salt)

	after := builder.String()
	after = MD5(after)

	after = string([]rune(after)[0:30])

	return after
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

// IsEmail 识别电子邮箱
func IsEmail(email string) bool {
	result, _ := regexp.MatchString(`^([\w\.\_\-]{2,10})@(\w{1,}).([a-z]{2,4})$`, email)

	return result
}

// CheckPasswordLever 密码强度必须为字⺟⼤⼩写+数字+符号，9位以上
func CheckPasswordLever(ps string) error {
	if len(ps) < 9 {
		return fmt.Errorf("密码至少需要9位")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("密码需要包含数字")
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("密码需要包含小写字符")
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
		return fmt.Errorf("密码需要包含大写字符")
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return fmt.Errorf("密码需要包含特殊字符")
	}
	return nil
}
