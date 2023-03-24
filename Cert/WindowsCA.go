package Cert

import (
	"encoding/pem"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
)

// 安装证书后设置Windows系统代理开启
// 参考blog：https://blog.csdn.net/leoforbest/article/details/120166881

// Person 个人
const Person = "my"

// TrustPeople 受信任人
const TrustPeople = "trustedpeople"

// TrustPublisher 可信任发行者
const TrustPublisher = "trustedpublishers"

// Root 受信任的根证书颁发机构
const Root = "root"

// Ca 证书颁发机构
const Ca = "ca"

// 证书安装
func Install() {
	cert, err := os.ReadFile("cert.crt")
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(cert)
	if block == nil {
		fmt.Println("解析证书内容失败：", err)
		return
	}
	certBytes := block.Bytes
	certContext, err := windows.CertCreateCertificateContext(windows.X509_ASN_ENCODING|windows.PKCS_7_ASN_ENCODING,
		&certBytes[0],
		uint32(len(certBytes)))
	if err != nil {
		fmt.Println("创建上下文失败：", err)
		return
	}
	defer func() {
		_ = windows.CertFreeCertificateContext(certContext)
	}()

	utf16Ptr, err := windows.UTF16PtrFromString(Root)
	storeHandle, err := windows.CertOpenSystemStore(0, utf16Ptr)
	if err != nil {
		fmt.Println("打开系统存储失败： ", err)
		return
	}
	defer func() {
		_ = windows.CertCloseStore(storeHandle, windows.CERT_CLOSE_STORE_FORCE_FLAG)
	}()

	err = windows.CertAddCertificateContextToStore(storeHandle, certContext, windows.CERT_STORE_ADD_REPLACE_EXISTING_INHERIT_PROPERTIES, nil)
	if err != nil {
		fmt.Println("安装失败 ", err)
		return
	}
	fmt.Println("安装成功")
}
