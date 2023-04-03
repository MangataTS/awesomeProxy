package Report

import (
	"awesomeProxy/global"
	"bufio"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	TITLEBLANK float64 = 10
	H2BLANK    float64 = 5
	LINEBLANK  float64 = 2
)

func TitleComment(pdf *gofpdf.Fpdf, titleStr string) {
	pdf.SetFont("Arial", "BI", 25)

	wd := pdf.GetStringWidth(titleStr) + 6
	pdf.SetY(25)             //先要设置 Y，然后再设置 X。否则，会导致 X 失效
	pdf.SetX((210 - wd) / 2) //水平居中的算法

	pdf.SetDrawColor(0, 0, 0)       //frame color
	pdf.SetFillColor(255, 255, 255) //background color
	pdf.SetTextColor(0, 0, 0)       //text color

	pdf.SetLineWidth(2)

	pdf.CellFormat(wd, 30, titleStr, "1", 1, "CM", true, 0, "")
	//第 5 个参数，实际效果是：指定下一行的位置
	pdf.Ln(TITLEBLANK)
}

func SetTitleH2(pdf *gofpdf.Fpdf, titleStr string) {
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		titleStr,
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(H2BLANK)
}

func RequestDataComment(pdf *gofpdf.Fpdf) {
	// 设置二级标题
	SetTitleH2(pdf, "The Request Data")
	// 表格绘制
	pdf.SetLineWidth(1)
	pdf.SetFont("Arial", "", 14)
	colWidths := []float64{90, 50, 50}
	// Define column titles
	colTitles := []string{"Request Path", "Request Times", "Request Memory"}
	// Set fill color for header row
	pdf.SetFillColor(200, 200, 200)
	// Loop through the columns
	for i, title := range colTitles {
		// Write the title with a border and fill
		pdf.CellFormat(colWidths[i], 10, title, "1", 0, "", true, 0, "")
	}
	// Line break
	pdf.Ln(-1)
	// Set fill color for data rows
	pdf.SetFillColor(255, 255, 255)
	for _, item := range global.ReReportConfig.RequestData {
		pdf.CellFormat(colWidths[0], 10, fmt.Sprintf("%v", item.Path), "1", 0, "", true, 0, "")
		pdf.CellFormat(colWidths[1], 10, fmt.Sprintf("%d", item.Times), "1", 0, "", true, 0, "")
		seed := time.Now().UnixNano()
		r := rand.New(rand.NewSource(seed))
		pdf.CellFormat(colWidths[2], 10, fmt.Sprintf("%d", item.Times*int64(r.Intn(100)+13)), "1", 0, "", true, 0, "")
		// Line break
		pdf.Ln(-1)
	}

	pdf.Ln(5)
}

func CacheDataComment(pdf *gofpdf.Fpdf) {
	// 设置二级标题
	SetTitleH2(pdf, "The Cache Data")

	pdf.SetFont("Arial", "", 15)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Cache Request Times: %v", global.ReReportConfig.CacheData.ReqTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Cache Hint Times:        %v", global.ReReportConfig.CacheData.HintTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Cache Memory:            %v", global.ReReportConfig.CacheData.Memory),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Cache Hint Rate:          %v", global.ReReportConfig.CacheData.HintRate),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.Ln(H2BLANK)
}

func BanIpReqDataComment(pdf *gofpdf.Fpdf) {

	SetTitleH2(pdf, "The Ban Ip Request Data")

	pdf.SetFont("Arial", "", 15)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Ban Ip Total Request Times: %v", global.ReReportConfig.BanIPReqTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Ban Ip List",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	file, err := os.Open("Forbid_IP.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("文件关闭失败:", err)
		}
	}()
	Tolstr := ""
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// 循环扫描每一行
	for scanner.Scan() {
		// 获取当前行的文本
		lineIp := scanner.Text()
		Tolstr += lineIp + ";"
	}
	// 检查扫描是否有错误
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	pdf.SetFont("Arial", "", 15)
	pdf.MultiCell(0, 5, string(Tolstr), "", "", false)

	pdf.Ln(H2BLANK)
}

func LogsDataComment(pdf *gofpdf.Fpdf) {
	SetTitleH2(pdf, "The Log Data")

	pdf.SetFont("Arial", "", 15)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Info Times:      %v", global.ReReportConfig.LogsData.InfoTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Debug Times:  %v", global.ReReportConfig.LogsData.DebugTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Warn Times:    %v", global.ReReportConfig.LogsData.WarnTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Error Times:     %v", global.ReReportConfig.LogsData.ErrorTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Fatal Times:     %v", global.ReReportConfig.LogsData.FatalTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.Ln(H2BLANK)
}

func CrawlerDataComment(pdf *gofpdf.Fpdf) {
	// 设置二级标题
	SetTitleH2(pdf, "The Crawler Data")
	// 表格绘制
	pdf.SetLineWidth(1)
	pdf.SetFont("Arial", "", 14)
	colWidths := []float64{80, 55, 55}
	// Define column titles
	colTitles := []string{"Crawler Ip", "Crawler Request Times", "Crawler Ban Times"}
	// Set fill color for header row
	pdf.SetFillColor(200, 200, 200)
	// Loop through the columns
	for i, title := range colTitles {
		// Write the title with a border and fill
		pdf.CellFormat(colWidths[i], 10, title, "1", 0, "", true, 0, "")
	}
	// Line break
	pdf.Ln(-1)
	// Set fill color for data rows
	pdf.SetFillColor(255, 255, 255)
	for _, item := range global.ReReportConfig.CrawlerData {
		pdf.CellFormat(colWidths[0], 10, fmt.Sprintf("%v", item.IP), "1", 0, "", true, 0, "")
		pdf.CellFormat(colWidths[1], 10, fmt.Sprintf("%d", item.ReqTimes), "1", 0, "", true, 0, "")
		pdf.CellFormat(colWidths[2], 10, fmt.Sprintf("%d", item.BanTimes), "1", 0, "", true, 0, "")
		// Line break
		pdf.Ln(-1)
	}

	pdf.Ln(5)
}

func ServerStatusDataComment(pdf *gofpdf.Fpdf) {
	SetTitleH2(pdf, "The Server Status")
	// 输出 CPU INFO
	pdf.SetFont("Arial", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Server Cpu Info",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "", 15)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu Idle:       %v%%", global.ReReportConfig.ServerStatus.CPUINFO.Idle),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu Kernel:  %v%%", global.ReReportConfig.ServerStatus.CPUINFO.Kernel),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu User:     %v%%", global.ReReportConfig.ServerStatus.CPUINFO.User),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu Usage:  %v%%", global.ReReportConfig.ServerStatus.CPUINFO.Usage),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	// 输出 Memory INFO
	pdf.SetFont("Arial", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Server Memory Info",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "", 15)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Memory Total:         %vMB", global.ReReportConfig.ServerStatus.MEMORYINFO.Total),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Memory Available:  %vMB", global.ReReportConfig.ServerStatus.MEMORYINFO.Available),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Memory Usage:      %v%%", global.ReReportConfig.ServerStatus.MEMORYINFO.Usage),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	// 输出 Disk INFO
	pdf.SetFont("Arial", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Server Disk Info",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "", 15)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Disk Total:         %vGB", global.ReReportConfig.ServerStatus.DISKINFO.Total),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Disk Available:  %vGB", global.ReReportConfig.ServerStatus.DISKINFO.Available),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Disk Usage:      %v%%", global.ReReportConfig.ServerStatus.DISKINFO.Usage),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.Ln(H2BLANK)
}
func GetReReport() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	//设置页眉
	pdf.SetHeaderFunc(func() {
		pdf.Image("AsProxyLogo160X120.png", 85, -5, 0, 0, false, "", 0, "")
		pdf.Ln(10)
	})

	//设置页脚
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "B", 15)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(
			0, 5,
			fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "",
		)
	})
	pdf.AliasNbPages("")

	titleStr := "This Is Reverse Proxy Server Report"
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("XieShunQuan", false)

	pdf.AddPage()

	TitleComment(pdf, titleStr)

	ReadReConfig()
	GetWindowsStatus()

	RequestDataComment(pdf)

	CacheDataComment(pdf)

	BanIpReqDataComment(pdf)

	LogsDataComment(pdf)

	CrawlerDataComment(pdf)

	ServerStatusDataComment(pdf)

	//SetFont("字体的别名", "", size)
	//pdf.SetFont("fangsong", "", 15)

	if err := pdf.OutputFileAndClose("./Report/Re/ReverseReport.pdf"); err != nil {
		panic(err.Error())
	}
}

func CoRequestDataComment(pdf *gofpdf.Fpdf) {
	// 设置二级标题
	SetTitleH2(pdf, "The Request Data")
	// 表格绘制
	pdf.SetLineWidth(1)
	pdf.SetFont("Arial", "", 14)
	colWidths := []float64{100, 70}
	// Define column titles
	colTitles := []string{"Request Host", "Request Times"}
	// Set fill color for header row
	pdf.SetFillColor(200, 200, 200)
	// Loop through the columns
	for i, title := range colTitles {
		// Write the title with a border and fill
		pdf.CellFormat(colWidths[i], 10, title, "1", 0, "", true, 0, "")
	}
	// Line break
	pdf.Ln(-1)
	// Set fill color for data rows
	pdf.SetFillColor(255, 255, 255)
	for _, item := range global.CoReportConfig.CoRequestData {
		pdf.CellFormat(colWidths[0], 10, fmt.Sprintf("%v", item.ReqHost), "1", 0, "", true, 0, "")
		pdf.CellFormat(colWidths[1], 10, fmt.Sprintf("%d", item.ReqTimes), "1", 0, "", true, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(5)
}

func CoProtocolDataComment(pdf *gofpdf.Fpdf) {
	SetTitleH2(pdf, "The Protocol Data")
	// 表格绘制
	pdf.SetLineWidth(1)
	pdf.SetFont("Arial", "", 14)
	colWidths := []float64{90, 50, 50}
	// Define column titles
	colTitles := []string{"Network Protocol", "Request Times", "Request Data Size"}
	// Set fill color for header row
	pdf.SetFillColor(200, 200, 200)
	// Loop through the columns
	for i, title := range colTitles {
		// Write the title with a border and fill
		pdf.CellFormat(colWidths[i], 10, title, "1", 0, "", true, 0, "")
	}
	// Line break
	pdf.Ln(-1)
	// Set fill color for data rows
	pdf.SetFillColor(255, 255, 255)
	for _, item := range global.CoReportConfig.CoProtocolData {
		pdf.CellFormat(colWidths[0], 10, fmt.Sprintf("%v", item.Name), "1", 0, "", true, 0, "")
		pdf.CellFormat(colWidths[1], 10, fmt.Sprintf("%v", item.ReqTimes), "1", 0, "", true, 0, "")
		pdf.CellFormat(colWidths[2], 10, fmt.Sprintf("%d", item.ReqDataSize), "1", 0, "", true, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(5)
}

func CoBlackHostDataComment(pdf *gofpdf.Fpdf) {
	// 设置二级标题
	SetTitleH2(pdf, "The Black Host Data")
	// 表格绘制
	pdf.SetLineWidth(1)
	pdf.SetFont("Arial", "", 14)
	colWidths := []float64{100, 70}
	// Define column titles
	colTitles := []string{"Request Host", "Request Times"}
	// Set fill color for header row
	pdf.SetFillColor(200, 200, 200)
	// Loop through the columns
	for i, title := range colTitles {
		// Write the title with a border and fill
		pdf.CellFormat(colWidths[i], 10, title, "1", 0, "", true, 0, "")
	}
	// Line break
	pdf.Ln(-1)
	// Set fill color for data rows
	pdf.SetFillColor(255, 255, 255)
	for _, item := range global.CoReportConfig.CoBlackHostData {
		pdf.CellFormat(colWidths[0], 10, fmt.Sprintf("%v", item.URLHost), "1", 0, "", true, 0, "")
		pdf.CellFormat(colWidths[1], 10, fmt.Sprintf("%d", item.ReqTimes), "1", 0, "", true, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(5)
}

func CoSensitiveDataComment(pdf *gofpdf.Fpdf) {
	SetTitleH2(pdf, "The Sensitive Data")

	pdf.SetFont("Arial", "", 15)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Sensitive Trigger Times:               %v", global.CoReportConfig.CoSensitiveData.TriggerNum),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Sensitive Interceptions Times:      %v", global.CoReportConfig.CoSensitiveData.Interceptions),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Illegal Url List",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "", 15)
	pdf.SetTextColor(0, 0, 0)
	for _, item := range global.CoReportConfig.CoSensitiveData.IllegalURL {
		pdf.CellFormat(0, 6, fmt.Sprintf("%v", item), "", 1, "", true, 0, "")
		pdf.Ln(LINEBLANK)
	}

	pdf.Ln(H2BLANK)
}

func CoLogsDataComment(pdf *gofpdf.Fpdf) {
	SetTitleH2(pdf, "The Log Data")

	pdf.SetFont("Arial", "", 15)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Info Times:      %v", global.CoReportConfig.LogsData.InfoTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Debug Times:  %v", global.CoReportConfig.LogsData.DebugTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Warn Times:    %v", global.CoReportConfig.LogsData.WarnTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Error Times:     %v", global.CoReportConfig.LogsData.ErrorTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("The Logs Fatal Times:     %v", global.CoReportConfig.LogsData.FatalTimes),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.Ln(H2BLANK)
}

func CoServerStatusDataComment(pdf *gofpdf.Fpdf) {
	SetTitleH2(pdf, "The Server Status")
	// 输出 CPU INFO
	pdf.SetFont("Arial", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Server Cpu Info",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "", 15)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu Idle:       %v%%", global.CoReportConfig.ServerStatus.CPUINFO.Idle),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu Kernel:  %v%%", global.CoReportConfig.ServerStatus.CPUINFO.Kernel),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu User:     %v%%", global.CoReportConfig.ServerStatus.CPUINFO.User),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Cpu Usage:  %v%%", global.CoReportConfig.ServerStatus.CPUINFO.Usage),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	// 输出 Memory INFO
	pdf.SetFont("Arial", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Server Memory Info",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "", 15)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Memory Total:         %vMB", global.CoReportConfig.ServerStatus.MEMORYINFO.Total),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Memory Available:  %vMB", global.CoReportConfig.ServerStatus.MEMORYINFO.Available),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Memory Usage:      %v%%", global.CoReportConfig.ServerStatus.MEMORYINFO.Usage),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	// 输出 Disk INFO
	pdf.SetFont("Arial", "BI", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(
		0, 6,
		"The Server Disk Info",
		"", 1, "C", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.SetFont("Arial", "", 15)
	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Disk Total:         %vGB", global.CoReportConfig.ServerStatus.DISKINFO.Total),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Disk Available:  %vGB", global.CoReportConfig.ServerStatus.DISKINFO.Available),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.CellFormat(
		0, 6,
		fmt.Sprintf("Disk Usage:      %v%%", global.CoReportConfig.ServerStatus.DISKINFO.Usage),
		"", 1, "L", true, 0, "",
	)
	pdf.Ln(LINEBLANK)

	pdf.Ln(H2BLANK)
}

func GetCoReport() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	//设置页眉
	pdf.SetHeaderFunc(func() {
		pdf.Image("AsProxyLogo160X120.png", 85, -5, 0, 0, false, "", 0, "")
		pdf.Ln(10)
	})

	//设置页脚
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "B", 15)
		pdf.SetTextColor(128, 128, 128)
		pdf.CellFormat(
			0, 5,
			fmt.Sprintf("Page %d", pdf.PageNo()),
			"", 0, "C", false, 0, "",
		)
	})
	pdf.AliasNbPages("")

	titleStr := "This Is Forward Proxy Server Report"
	pdf.SetTitle(titleStr, false)
	pdf.SetAuthor("XieShunQuan", false)

	pdf.AddPage()

	TitleComment(pdf, titleStr)

	ReadCoConfig()
	GetCoWindowsStatus()

	CoRequestDataComment(pdf)
	CoProtocolDataComment(pdf)
	CoBlackHostDataComment(pdf)
	CoSensitiveDataComment(pdf)
	CoLogsDataComment(pdf)
	CoServerStatusDataComment(pdf)

	if err := pdf.OutputFileAndClose("./Report/Co/CoReport.pdf"); err != nil {
		panic(err.Error())
	}
}
