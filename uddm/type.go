package uddm

// Mobile 手机号 132****7986
type Mobile string

// BankCard 银行卡号 622888******5676
type BankCard string

// IDCard 身份证号 12525252******5252
type IDCard string

// IDName 姓名 *鸿章
// TODO: 有更好的性能选择 https://blog.thinkeridea.com/201910/go/efficient_string_truncation.html
type IDName string

// PassWord 密码 ******
type PassWord string

// Email 邮箱 l***w@gmail.com
type Email string
