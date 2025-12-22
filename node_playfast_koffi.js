const koffi = require('koffi');
const path = require('path');
const os = require('os');

// 根据平台和架构选择库文件
function getLibraryPath() {
  const platform = os.platform();
  const arch = os.arch(); // 'x64' 或 'ia32'

  let libName;
  if (platform === 'win32') {
    libName = arch === 'x64' ? 'signBox64.dll' : 'signBox32.dll';
  } else if (platform === 'darwin') {
    libName = arch === 'x64' ? 'signBox64.dylib' : 'signBox32.dylib';
  } else {
    libName = arch === 'x64' ? 'signBox64.so' : 'signBox32.so';
  }
  
  return path.join(__dirname, libName);
}

// 加载共享库
const libPath = getLibraryPath();
console.log('加载库文件:', libPath);

const lib = koffi.load(libPath);

// 定义函数签名
const Init = lib.func('Init', 'void', []);
const Switch = lib.func('Switch', 'str', ['int', 'str', 'int','str']);
const GetProxyList = lib.func('GetProxyList', 'str', []);
const GetAnnouncement = lib.func('GetAnnouncement', 'str', []);
const GetVersion = lib.func('GetVersion', 'str', []);
const GetRouteInfo = lib.func('GetRouteInfo', 'str', []);
const Stop = lib.func('Stop', 'void', []);
const Cleanup = lib.func('Cleanup', 'void', []);
const FreeString = lib.func('FreeString', 'void', ['str']);

// 辅助函数：安全释放字符串
// 注意：koffi 返回的字符串已经是 JavaScript 字符串副本，理论上不需要释放
// 但如果 DLL 返回的是 C 字符串指针，可能需要释放
// 这里先尝试复制字符串内容，然后释放原始指针
function safeFree(str) {
  // koffi 的 'str' 类型已经自动处理了内存，通常不需要手动释放
  // 如果确实需要释放（根据你的 DLL 实现），可以取消下面的注释
  // try {
  //   if (str && typeof str === 'string') {
  //     // koffi 返回的字符串是副本，不需要释放
  //   }
  // } catch (error) {
  //   console.error('释放字符串失败:', error);
  // }
}

// 主函数
async function main() {
  try {
    // console.log('=== PlayFast 共享库调用示例 ===\n');

    // 1. 初始化
    console.log('[1] 初始化...');
    Init();
    console.log('✓ 初始化成功\n');

    // 2. 获取版本号
    // console.log('[2] 获取版本号...');
    // try {
    //   const version = GetVersion();
    //   console.log('✓ 版本:', version);
    //   console.log('版本类型:', typeof version);
    //   console.log('版本长度:', version ? version.length : 0);
    //   // koffi 返回的字符串不需要手动释放
    //   // safeFree(version);
    //   console.log('');
    // } catch (error) {
    //   console.error('✗ 获取版本号失败:', error);
    //   console.error('错误堆栈:', error.stack);
    //   throw error;
    // }

    // 3. 获取代理列表
    console.log('[3] 获取代理列表...');
    let proxyList = [];
    try {
      const proxyListJson = GetProxyList();
      console.log('代理列表原始数据:', proxyListJson);
      console.log('代理列表类型:', typeof proxyListJson);
      try {
        proxyList = JSON.parse(proxyListJson);
        console.log('✓ 可用代理节点:', proxyList);
      } catch (parseError) {
        console.error('✗ 解析代理列表失败:', parseError);
        console.log('原始数据:', proxyListJson);
      }
      // koffi 返回的字符串不需要手动释放
      // safeFree(proxyListJson);
      console.log('');
      
      if (proxyList.length === 0) {
        console.log('没有可用的代理节点，退出');
        Cleanup();
        return;
      }
    } catch (error) {
      console.error('✗ 获取代理列表失败:', error);
      console.error('错误堆栈:', error.stack);
      Cleanup();
      return;
    }

    // 4. 启动加速（使用第一个节点）
    console.log('[4] 启动加速...');
    try {
      const proxyName = proxyList[0];
      console.log('  使用节点:', proxyName);


      //指定程序加速
      const appList = JSON.stringify(['chrome.exe','msedge.exe']);
      const errorMsg = Switch(1,"香港(vless)", 0,appList); // 1=启动, 0=不启用路由模式
      
      if (errorMsg && errorMsg.length > 0) {
        console.error('✗ 启动失败:', errorMsg);
        // safeFree(errorMsg);
        Cleanup();
        return;
      }
      
      console.log('✓ 加速已启动');
      // safeFree(errorMsg);
      console.log('');
    } catch (error) {
      console.error('✗ 启动加速失败:', error);
      console.error('错误堆栈:', error.stack);
      Cleanup();
      return;
    }

    // 5. 等待一段时间
    console.log('[5] 运行中...');
    //60000 60
    await new Promise(resolve => setTimeout(resolve, (60000) * 1));

    // 6. 停止加速
    console.log('[6] 停止加速...');
    Stop();
    console.log('✓ 加速已停止\n');



    // 8. 清理资源
    console.log('[8] 清理资源...');
    Cleanup();
    console.log('✓ 清理完成\n');



  } catch (error) {
    console.error('运行出错:', error);
    console.error('错误堆栈:', error.stack);
    try {
      Cleanup();
    } catch (cleanupError) {
      console.error('清理资源时出错:', cleanupError);
    }
  }
}

// 运行
main();

