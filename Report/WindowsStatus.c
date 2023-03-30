//#include <windows.h>
//#include <stdio.h>
//
//// convert filetime to 64-bit integer
//__int64 filetimeToInt64(const FILETIME* ft) {
//	return (((__int64)ft->dwHighDateTime) << 32) | ft->dwLowDateTime;
//}
//
//// convert filetime to double
//double filetimeToDouble(const FILETIME* ft) {
//	return (double)filetimeToInt64(ft);
//}
//
//// convert filetime to seconds
//double filetimeToSeconds(const FILETIME* ft) {
//	return filetimeToDouble(ft) / 10000000.0;
//}
//
//struct CPU_INFO {
//	double idle;
//	double kernel;
//	double user;
//	double Usage;
//};
//
//double GetCpuUsage() {
//	FILETIME idleTime1, kernelTime1, userTime1; // 第一次调用GetSystemTimes的结果
//	FILETIME idleTime2, kernelTime2, userTime2; // 第二次调用GetSystemTimes的结果
//	ULARGE_INTEGER idle1, kernel1, user1; // 用于存储64位整数
//	ULARGE_INTEGER idle2, kernel2, user2;
//	ULONGLONG idleDiff, kernelDiff, userDiff; // 用于存储时间差
//	ULONGLONG totalDiff; // 用于存储总时间差
//	double cpuUsage; // 用于存储CPU利用率
//
//	// 第一次调用GetSystemTimes，获取系统启动后的内核模式时间、用户模式时间和空闲时间
//	if (!GetSystemTimes(&idleTime1, &kernelTime1, &userTime1)) {
//		printf("GetSystemTimes failed\n");
//		return -1;
//	}
//
//	// 等待一秒钟
//	Sleep(1000);
//
//	// 第二次调用GetSystemTimes，获取系统启动后的内核模式时间、用户模式时间和空闲时间
//	if (!GetSystemTimes(&idleTime2, &kernelTime2, &userTime2)) {
//		printf("GetSystemTimes failed\n");
//		return -1;
//	}
//
//	// 将FILETIME结构转换为64位整数
//	idle1.LowPart = idleTime1.dwLowDateTime;
//	idle1.HighPart = idleTime1.dwHighDateTime;
//	kernel1.LowPart = kernelTime1.dwLowDateTime;
//	kernel1.HighPart = kernelTime1.dwHighDateTime;
//	user1.LowPart = userTime1.dwLowDateTime;
//	user1.HighPart = userTime1.dwHighDateTime;
//
//	idle2.LowPart = idleTime2.dwLowDateTime;
//	idle2.HighPart = idleTime2.dwHighDateTime;
//	kernel2.LowPart = kernelTime2.dwLowDateTime;
//	kernel2.HighPart = kernelTime2.dwHighDateTime;
//	user2.LowPart = userTime2.dwLowDateTime;
//	user2.HighPart = userTime2.dwHighDateTime;
//
//	// 计算两次调用之间的时间差（单位是100纳秒）
//	idleDiff = idle2.QuadPart - idle1.QuadPart;
//	kernelDiff = kernel2.QuadPart - kernel1.QuadPart;
//	userDiff = user2.QuadPart - user1.QuadPart;
//
//	// 计算总时间差（单位是100纳秒）
//	totalDiff = kernelDiff + userDiff;
//
//	// 计算CPU利用率（百分比）
//	cpuUsage = 100.0 * (totalDiff - idleDiff) / totalDiff;
//
//	// 打印CPU利用率
//	//printf("CPU usage: %.2f%%\n", cpuUsage);
//	return cpuUsage;
//
//}
//
//struct CPU_INFO GetCpuInfo() {
//	// get cpu usage
//	FILETIME idleTime, kernelTime, userTime;
//	struct CPU_INFO ci;
//	BOOL result = GetSystemTimes(&idleTime, &kernelTime, &userTime);
//	if (!result) {
//		printf("GetSystemTimes failed: %d\n", GetLastError());
//		return  ci;
//	}
//	double idleSeconds = filetimeToSeconds(&idleTime);
//	double kernelSeconds = filetimeToSeconds(&kernelTime);
//	double userSeconds = filetimeToSeconds(&userTime);
//	double totalSeconds = kernelSeconds + userSeconds;
//	ci.idle = idleSeconds / totalSeconds * 100;//CPU处于空闲状态时间比例
//	ci.kernel = (kernelSeconds - idleSeconds) / totalSeconds * 100;
//	ci.user = userSeconds / totalSeconds * 100;
//	ci.Usage = GetCpuUsage();
//	//	printf("CPU idle: %.2f%%\n", idleSeconds / totalSeconds * 100);
//	//	printf("CPU kernel: %.2f%%\n", (kernelSeconds - idleSeconds) / totalSeconds * 100);//CPU处于系统内核执行的时间比例
//	//	printf("CPU user: %.2f%%\n", userSeconds / totalSeconds * 100);//CPU处于用户态执行的时间比例
//	return ci;
//}
//
//struct MemAndDiskInfo{
//	unsigned long  long total;		//Memory=MB,Disk=GB
//	unsigned long  long available;	//Memory=MB,Disk=GB
//	double 				usage;		//使用率
//};
//
//struct MemAndDiskInfo GetMemoryInfo() {
//	// get memory status
//	MEMORYSTATUSEX memStatus;
//	struct MemAndDiskInfo Mi;
//	memStatus.dwLength = sizeof(memStatus);
//	BOOL result = GlobalMemoryStatusEx(&memStatus);
//	if (!result) {
//		printf("GlobalMemoryStatusEx failed: %d\n", GetLastError());
//		return Mi;
//	}
//	Mi.total = memStatus.ullTotalPhys / 1024 / 1024;
//	Mi.available = memStatus.ullAvailPhys / 1024 / 1024;
//	Mi.usage = (double)(memStatus.ullTotalPhys - memStatus.ullAvailPhys) / memStatus.ullTotalPhys * 100;
//	//	printf("Memory total: %llu MB\n", memStatus.ullTotalPhys / 1024 / 1024);
//	//	printf("Memory available: %llu MB\n", memStatus.ullAvailPhys / 1024 / 1024);
//	//	printf("Memory usage: %.2f%%\n", (double)(memStatus.ullTotalPhys - memStatus.ullAvailPhys) / memStatus.ullTotalPhys * 100);
//	return Mi;
//}
//
//
//struct MemAndDiskInfo GetDiskInfo() {
//	// get disk free space
//	ULARGE_INTEGER freeBytes, totalBytes, totalFreeBytes;
//	struct MemAndDiskInfo Di;
//	BOOL result = GetDiskFreeSpaceEx(NULL, &freeBytes, &totalBytes, &totalFreeBytes);
//	if (!result) {
//		printf("GetDiskFreeSpaceEx failed: %d\n", GetLastError());
//		return Di;
//	}
//	Di.total = totalBytes.QuadPart / 1024 / 1024 / 1024;
//	Di.available = freeBytes.QuadPart / 1024 / 1024 / 1024;
//	Di.usage = (double)(totalBytes.QuadPart - freeBytes.QuadPart) / totalBytes.QuadPart * 100;
//	//	printf("Disk total: %llu GB\n", totalBytes.QuadPart / 1024 / 1024 / 1024);
//	//	printf("Disk available: %llu GB\n", freeBytes.QuadPart / 1024 / 1024 / 1024);
//	//	printf("Disk usage: %.2f%%\n", (double)(totalBytes.QuadPart - freeBytes.QuadPart) / totalBytes.QuadPart * 100);
//	return Di;
//}