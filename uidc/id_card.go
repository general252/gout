package uidc

import (
	"fmt"
	"strconv"
)

// IdCardNumberMask 身份证敏感处理
func IdCardNumberMask(idCardNumber string) string {
	if len(idCardNumber) >= 14 {
		return idCardNumber[:10] + "****" + idCardNumber[14:]
	} else {
		return "****"
	}
}

// GetBirthdayFromIdCardNumber 根据身份证获取生日
func GetBirthdayFromIdCardNumber(idCardNumber string) (int, int, int, error) {
	if IdCardNumberCheck(idCardNumber) == false {
		return 0, 0, 0, fmt.Errorf("error id card number")
	}

	birthday := idCardNumber[6:14]
	year, err1 := strconv.Atoi(birthday[0:4])
	month, err2 := strconv.Atoi(birthday[4:6])
	day, err3 := strconv.Atoi(birthday[6:8])
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, fmt.Errorf("error id card number")
	}

	return year, month, day, nil
}

// IdCardNumberCheck 验证身份证是否有效
func IdCardNumberCheck(idCardNumber string) bool {
	/**
	  中国居民身份证号码编码规则
	  第一、二位表示省（自治区、直辖市、特别行政区）。
	  第三、四位表示市（地级市、自治州、盟及国家直辖市所属市辖区和县的汇总码）。其中，01-20，51-70表示省直辖市；21-50表示地区（自治州、盟）。
	  第五、六位表示县（市辖区、县级市、旗）。01-18表示市辖区或地区（自治州、盟）辖县级市；21-80表示县（旗）；81-99表示省直辖县级市。
	  第七、十四位表示出生年月日（单数字月日左侧用0补齐）。其中年份用四位数字表示，年、月、日之间不用分隔符。例如：1981年05月11日就用19810511表示。
	  第十五、十七位表示顺序码。对同地区、同年、月、日出生的人员编定的顺序号。其中第十七位奇数分给男性，偶数分给女性。
	  第十八位表示校验码。作为尾号的校验码，是由号码编制单位按统一的公式计算出来的，校验码如果出现数字10，就用X来代替，详情参考下方计算方法。

	  其中第一代身份证号码为15位。年份两位数字表示，没有校验码。
	  前六位详情请参考省市县地区代码
	  X是罗马字符表示数字10
	*/

	/**
	中国居民身份证校验码算法
	步骤如下：

	1. 将身份证号码前面的17位数分别乘以不同的系数。从第一位到第十七位的系数分别为：7－9－10－5－8－4－2－1－6－3－7－9－10－5－8－4－2。
	2. 将这17位数字和系数相乘的结果相加。
	3. 用加出来和除以11，取余数。
	4. 余数只可能有0－1－2－3－4－5－6－7－8－9－10这11个数字。其分别对应的最后一位身份证的号码为1－0－X－9－8－7－6－5－4－3－2。
	5. 通过上面计算得知如果余数是3，第18位的校验码就是9。如果余数是2那么对应的校验码就是X，X实际是罗马数字10。
	例如
		某男性的身份证号码为【53010219200508011x】， 我们看看这个身份证是不是合法的身份证。
	首先我们得出前17位的乘积和
	【(5*7)+(3*9)+(0*10)+(1*5)+(0*8)+(2*4)+(1*2)+(9*1)+(2*6)+(0*3)+(0*7)+(5*9)+(0*10)+(8*5)+(0*8)+(1*4)+(1*2)】
	是189，然后用189除以11得出的结果是189/11=17----2，
	也就是说其余数是2。最后通过对应规则就可以知道余数2对应的检验码是X。所以，可以判定这是一个正确的身份证号码。
	*/

	// 中国居民身份证号码编码规则 http://www.ip33.com/shenfenzheng.html
	// 身份证生成 http://sfz.uzuzuz.com/
	// 中国省市县地区代码 http://www.ip33.com/area/index.html

	if len(idCardNumber) != 18 {
		return false
	}

	var bits = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'} // 11
	var numbers = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2} // 17

	var sum = 0
	for i := 0; i < len(idCardNumber)-1; i++ {
		sum += int(idCardNumber[i]-'0') * numbers[i]
	}

	if bits[sum%len(bits)] == idCardNumber[17] {
		return true
	}

	return false
}
