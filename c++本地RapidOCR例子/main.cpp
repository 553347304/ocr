// version 2017  需要引入opencv
// json头文件 https://github.com/nlohmann/json/releases/download/v3.10.5/json.hpp
#include <Windows.h>
#include <iostream>
#include <vector>
#include <string>
#include <regex>
#include <Shlwapi.h>
#include <opencv2/opencv.hpp>
#pragma comment(lib, "opencv_world490.lib")
#pragma comment(lib, "Shlwapi.lib")
using namespace std;

#include <cstdio>
#include <memory>
#include <stdexcept>
#include <array>
#include "json.hpp"
using json = nlohmann::json;
namespace fs = std::filesystem;

struct FIND { int X; int Y; };
FIND DPI;
cv::Mat IMG;
cv::Mat img;

string WCHAR_STR(const wchar_t* wc) {
    int len = WideCharToMultiByte(CP_ACP, 0, wc, int(wcslen(wc)), NULL, 0, NULL, NULL);
    char* m_char = new char[len + 1];
    WideCharToMultiByte(CP_ACP, 0, wc, int(wcslen(wc)), m_char, len, NULL, NULL);
    m_char[len] = '\0';
    string result(m_char);
    delete[] m_char;
    return result;
}
const wchar_t* STR_WCHAR(string value) {
    // string_wstring
    LPCSTR pszSrc = value.c_str();
    int nLen = MultiByteToWideChar(CP_ACP, 0, pszSrc, -1, NULL, 0);
    if (nLen == 0) return L"";

    wchar_t* pwszDst = new wchar_t[nLen];
    if (!pwszDst) return L"";

    MultiByteToWideChar(CP_ACP, 0, pszSrc, -1, pwszDst, nLen);
    wstring wstr(pwszDst);
    delete[] pwszDst;
    pwszDst = NULL;

    // wstring_const wchar_t*
    const wchar_t* pwidstr = wstr.c_str();
    return pwidstr;
}
cv::Mat HBITMAP_Mat(HBITMAP hBitmap) {
    BITMAP bmp;
    GetObject(hBitmap, sizeof(BITMAP), &bmp);
    HDC hdc = GetDC(NULL);
    HDC memDC = CreateCompatibleDC(hdc);
    SelectObject(memDC, hBitmap);
    BITMAPINFOHEADER bi = { sizeof(BITMAPINFOHEADER), bmp.bmWidth, bmp.bmHeight, 1, 32, BI_RGB, 0, 0, 0, 0 };
    cv::Mat image(bmp.bmHeight, bmp.bmWidth, CV_8UC4);
    GetDIBits(hdc, hBitmap, 0, bmp.bmHeight, image.data, (BITMAPINFO*)&bi, DIB_RGB_COLORS);
    cv::Mat mat;
    cv::flip(image, mat, 0);
    DeleteDC(memDC);
    ReleaseDC(NULL, hdc);
    return mat;
};

string Utf8ToSystem(const string& utf8_string) {
    int len = MultiByteToWideChar(CP_UTF8, 0, utf8_string.c_str(), -1, NULL, 0);
    vector<wchar_t> wide_string(len);
    MultiByteToWideChar(CP_UTF8, 0, utf8_string.c_str(), -1, &wide_string[0], len);

    len = WideCharToMultiByte(GetConsoleOutputCP(), 0, &wide_string[0], -1, NULL, 0, NULL, NULL);
    vector<char> system_string(len);
    WideCharToMultiByte(GetConsoleOutputCP(), 0, &wide_string[0], -1, &system_string[0], len, NULL, NULL);

    return string(&system_string[0]);
}

string ExecCommand(const string& command) {
    array<char, 128> buffer;
    string result;

    SECURITY_ATTRIBUTES saAttr;
    saAttr.nLength = sizeof(SECURITY_ATTRIBUTES);
    saAttr.bInheritHandle = TRUE;
    saAttr.lpSecurityDescriptor = NULL;

    HANDLE hReadPipe, hWritePipe;
    if (!CreatePipe(&hReadPipe, &hWritePipe, &saAttr, 0)) {
        throw std::runtime_error("CreatePipe() failed!");
    }

    STARTUPINFOA si = { sizeof(si) };
    si.dwFlags |= STARTF_USESHOWWINDOW | STARTF_USESTDHANDLES;
    si.wShowWindow = SW_HIDE; // 隐藏窗口
    si.hStdOutput = hWritePipe;
    si.hStdError = hWritePipe;

    PROCESS_INFORMATION pi;
    if (!CreateProcessA(NULL, const_cast<char*>(command.c_str()), NULL, NULL, TRUE, CREATE_NO_WINDOW, NULL, NULL, &si, &pi)) {
        CloseHandle(hReadPipe);
        CloseHandle(hWritePipe);
        throw std::runtime_error("CreateProcessA() failed!");
    }

    CloseHandle(hWritePipe); // 关闭写端，只保留读端

    DWORD bytesRead;
    while (ReadFile(hReadPipe, buffer.data(), buffer.size(), &bytesRead, NULL) && bytesRead != 0) {
        result.append(buffer.data(), bytesRead);
    }

    CloseHandle(hReadPipe);

    WaitForSingleObject(pi.hProcess, INFINITE);
    DWORD returnCode;
    GetExitCodeProcess(pi.hProcess, &returnCode);

    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);

    if (returnCode != 0) {
        std::cerr << "Command executed with return code: " << returnCode << endl;
    }
    return result;
}

struct RapidocrResponseBox {
    vector<vector<int>> box;
    double score;
    string text;
};
struct RapidocrResponse {
    int code;
    vector<RapidocrResponseBox> data;
    friend ostream& operator<<(ostream& os, const RapidocrResponse& r) {
        os << "code: " << r.code << endl;
        for (const auto& data : r.data) {
            os << "box: [";
            for (const auto& c : data.box) {
                os << "[" << c[0] << "," << c[1] << "]";
            }
            os << "] ";
            os << "Score: " << data.score << " text: " << data.text << endl;
        }
        return os;
    }
};
struct RapidocrResult {
    bool state;
    vector<string> text;
    vector<FIND> data;
    friend ostream& operator<<(ostream& os, const RapidocrResult& r) {
        os << "state: " << boolalpha << r.state << " ";
        int length = end(r.data) - begin(r.data);
        for (size_t i = 0; i < length; i++) {
            os << "{ ";
            os << "x: " << r.data[i].X << ", ";
            os << "y: " << r.data[i].Y;
            os << " } ";
            os << "text: " << r.text[i] << endl;
        }
        return os;
    }
};
class BAIYIN {
private:
    string ImagePath;
    string Extension;
    string _Pic;
    cv::Mat _mask;

    RapidocrResponse Response(cv::Mat img) {
        // 获取自身目录
        wchar_t selfPath[MAX_PATH];
        if (GetModuleFileName(nullptr, selfPath, MAX_PATH) > 0) PathRemoveFileSpec(selfPath);

        string dir = WCHAR_STR(selfPath) + "\\RapidOCR";

        string model = "v4";
        string imgPath = dir + "\\temp.png";
        cv::imwrite(imgPath, img);  // 保存临时图片

        // cmd
        string command = dir + "\\RapidOCR-json.exe"
            + " --models=models"
            + " --det=ch_PP-OCR" + model + "_det_infer.onnx"
            + " --cls=ch_ppocr_mobile_v2.0_cls_infer.onnx"
            + " --rec=ch_PP-OCR" + model + "_rec_infer.onnx"
            + " --keys=ppocr_keys_v1.txt"
            + " --image_path=\"" + imgPath + "\"";
        SetCurrentDirectory(STR_WCHAR(dir));    // 设置工作目录

        string exec = ExecCommand(command);     // 执行命令

        // 正则提取json
        regex pattern(R"(\{.*\})");
        smatch match;
        string regexp = "";
        if (regex_search(exec, match, pattern)) regexp = match[0];




        // 解析到结构体
        json j = json::parse(regexp);
        RapidocrResponse response;
        response.code = j["code"];
        for (const auto& item : j["data"]) {
            RapidocrResponseBox data;
            data.box = item["box"].get<vector<vector<int>>>();
            data.score = item["score"];
            data.text = item["text"];
            data.text = Utf8ToSystem(data.text);
            response.data.push_back(data);
        }
        return response;
    }
public:
    BAIYIN() {}
    FIND Window_Border(const char* className, const char* windowName) {

        HWND hwnd = FindWindowA(className, windowName);  // 获取窗口句柄
        if (hwnd == NULL) { return { -1, -1 }; }

        // 还原窗口>>>置顶窗口>>>移动缩放窗口
        ShowWindow(hwnd, SW_RESTORE);
        SetForegroundWindow(hwnd);
        SetWindowPos(hwnd, NULL, 0, 0, 1280, 720, SWP_NOZORDER | SWP_NOSIZE | SWP_SHOWWINDOW);

        RECT windowRect, clientRect;
        GetWindowRect(hwnd, &windowRect);  // 获取窗口大小
        GetClientRect(hwnd, &clientRect);  // 获取客户端大小

        // 边框间距
        int windowSizeX = windowRect.right - clientRect.right;
        int windowSizeY = windowRect.bottom - clientRect.bottom;

        // 边框大小
        int borderX = windowRect.left + windowSizeX / 2;
        int borderY = windowRect.top + windowSizeY - windowSizeX / 2;

        Sleep(100);

        // 移动缩放窗口
        MoveWindow(hwnd, 0, 0, 1280 + windowSizeX, 720 + windowSizeY, TRUE);

        return { borderX, borderY };
    }
    cv::Mat CaptureScreen(HWND hwnd = NULL, int x = 0, int y = 0, int width = 1280, int height = 720) {
        x = x + DPI.X;  y = y + DPI.Y;

        // 获取句柄窗口大小
        if (hwnd != NULL) {
            RECT rc;
            GetClientRect(hwnd, &rc);
            width = rc.right - rc.left;
            height = rc.bottom - rc.top;
        }

        // 创建屏幕设备上下文
        HDC hWindowDC = GetDC(hwnd);
        HDC hMemDC = CreateCompatibleDC(hWindowDC);
        HBITMAP hBitmap = CreateCompatibleBitmap(hWindowDC, width, height);
        HGDIOBJ oldBitmap = SelectObject(hMemDC, hBitmap);

        // 保存截图>>>内存
        BitBlt(hMemDC, 0, 0, width, height, hWindowDC, x, y, SRCCOPY);
        // 翻转图像
        // StretchBlt(hMemDC, 0, height - 1, width, -height, hWindowDC, x, y, width, height, SRCCOPY);

        // 将HBITMAP转换为CV图>>>灰度图>>>RGB
        cv::Mat IMG = HBITMAP_Mat(hBitmap);

        try {
            cv::cvtColor(IMG, IMG, cv::COLOR_BGRA2RGBA);
            cv::cvtColor(IMG, IMG, cv::COLOR_BGRA2RGB);
        }
        catch (const std::exception& e) {
            cout << "截图失败：跳过截图" << IMG.empty() << endl;
            std::cerr << "截图失败：" << e.what() << std::endl;
            //IMG_Bool = false;
            return CaptureScreen();  // 重新执行截图
        }

        //Save_Pic(IMG, IMG);  // 保存图片>>>本地

        // 释放资源
        SelectObject(hMemDC, oldBitmap);
        DeleteObject(hBitmap);
        DeleteDC(hMemDC);
        ReleaseDC(hwnd, hWindowDC);

        return IMG;
    }
    RapidocrResult ocrRapid(cv::Rect zoom = { 0, 0, 1280, 720 }) {
        // 裁剪图片>>>灰度图
        if (zoom.empty()) { zoom = { 0, 0, 1280, 720 }; }
        zoom = { zoom.x, zoom.y, zoom.width - zoom.x, zoom.height - zoom.y };
        cv::cvtColor(IMG, IMG, cv::COLOR_BGR2GRAY);
        cv::Mat img = IMG(zoom);

        RapidocrResult result;
        result.state = false;
        try {
            RapidocrResponse response = this->Response(img);
            for (const auto& data : response.data) {
                FIND find;
                find.X = zoom.x + data.box[0][0] + (data.box[2][0] - data.box[0][0]) / 2;
                find.Y = zoom.y + data.box[0][1] + (data.box[2][1] - data.box[0][1]) / 2;
                result.data.push_back(find);
                result.text.push_back(data.text);
                result.state = true;
            }
        }
        catch (const std::exception&) {
            result.state = false;
        }

        return result;
    }
};
BAIYIN BY;

int main() {
    //DPI = BY.Window_Border("FolderView", "SysListView32");
    IMG = BY.CaptureScreen();
    RapidocrResult result = BY.ocrRapid({ 0,0,1280,720 });

    cout << result << endl;
    return 0;
}