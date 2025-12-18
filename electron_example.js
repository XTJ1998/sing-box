/**
 * Electron 使用 PlayFast DLL/dylib 的完整示例
 * 
 * 需要安装依赖：
 * npm install ffi-napi ref-napi
 */

const ffi = require('ffi-napi');
const ref = require('ref-napi');
const path = require('path');
const os = require('os');

// 根据平台和架构选择库文件
function getLibraryPath() {
  const platform = os.platform();
  const arch = os.arch();
  
  let libName;
  if (platform === 'win32') {
    libName = arch === 'x64' ? 'playfast_amd64.dll' : 'playfast_386.dll';
  } else if (platform === 'darwin') {
    libName = arch === 'x64' ? 'playfast_amd64.dylib' : 'playfast_386.dylib';
  } else {
    libName = arch === 'x64' ? 'playfast_amd64.so' : 'playfast_386.so';
  }
  
  return path.join(__dirname, libName);
}

// 定义 C 函数接口
const playfast = ffi.Library(getLibraryPath(), {
  'Init': ['void', []],
  'Switch': ['string', ['int', 'string', 'int']],
  'GetProxyList': ['string', []],
  'GetAnnouncement': ['string', []],
  'GetVersion': ['string', []],
  'GetRouteInfo': ['string', []],
  'Stop': ['void', []],
  'Cleanup': ['void', []],
  'FreeString': ['void', ['string']]
});

class PlayFast {
  constructor() {
    this.initialized = false;
  }

  /**
   * 初始化核心模块
   */
  init() {
    if (this.initialized) {
      console.warn('PlayFast 已经初始化');
      return;
    }
    try {
      playfast.Init();
      this.initialized = true;
      console.log('PlayFast 初始化成功');
    } catch (error) {
      console.error('初始化失败:', error);
      throw error;
    }
  }

  /**
   * 获取版本号
   */
  getVersion() {
    if (!this.initialized) {
      throw new Error('请先调用 init()');
    }
    const version = playfast.GetVersion();
    const result = version;
    playfast.FreeString(version);
    return result;
  }

  /**
   * 获取代理列表
   */
  getProxyList() {
    if (!this.initialized) {
      throw new Error('请先调用 init()');
    }
    const jsonStr = playfast.GetProxyList();
    try {
      const list = JSON.parse(jsonStr);
      playfast.FreeString(jsonStr);
      return list;
    } catch (error) {
      playfast.FreeString(jsonStr);
      throw new Error('解析代理列表失败: ' + error.message);
    }
  }

  /**
   * 获取公告
   */
  getAnnouncement() {
    if (!this.initialized) {
      throw new Error('请先调用 init()');
    }
    const announcement = playfast.GetAnnouncement();
    const result = announcement;
    playfast.FreeString(announcement);
    return result;
  }

  /**
   * 启动或停止加速
   * @param {boolean} status - true=启动, false=停止
   * @param {string} proxy - 代理节点名称
   * @param {boolean} route - 是否启用路由模式
   * @returns {string} 错误信息，成功返回空字符串
   */
  switch(status, proxy, route = false) {
    if (!this.initialized) {
      throw new Error('请先调用 init()');
    }
    const statusInt = status ? 1 : 0;
    const routeInt = route ? 1 : 0;
    const error = playfast.Switch(statusInt, proxy, routeInt);
    const errorMsg = error;
    playfast.FreeString(error);
    return errorMsg;
  }

  /**
   * 获取路由配置信息
   */
  getRouteInfo() {
    if (!this.initialized) {
      throw new Error('请先调用 init()');
    }
    const jsonStr = playfast.GetRouteInfo();
    try {
      const info = JSON.parse(jsonStr);
      playfast.FreeString(jsonStr);
      return info;
    } catch (error) {
      playfast.FreeString(jsonStr);
      throw new Error('解析路由信息失败: ' + error.message);
    }
  }

  /**
   * 停止加速
   */
  stop() {
    if (!this.initialized) {
      return;
    }
    try {
      playfast.Stop();
      console.log('加速已停止');
    } catch (error) {
      console.error('停止加速失败:', error);
    }
  }

  /**
   * 清理资源
   */
  cleanup() {
    if (!this.initialized) {
      return;
    }
    try {
      this.stop();
      playfast.Cleanup();
      this.initialized = false;
      console.log('资源已清理');
    } catch (error) {
      console.error('清理资源失败:', error);
    }
  }
}

// 使用示例
async function example() {
  const pf = new PlayFast();
  
  try {
    // 初始化
    pf.init();
    
    // 获取版本
    console.log('版本:', pf.getVersion());
    
    // 获取代理列表
    const proxies = pf.getProxyList();
    console.log('可用代理:', proxies);
    
    if (proxies.length > 0) {
      // 启动加速（不使用路由模式）
      const error = pf.switch(true, proxies[0], false);
      if (error) {
        console.error('启动失败:', error);
      } else {
        console.log('加速已启动');
        
        // 等待一段时间...
        await new Promise(resolve => setTimeout(resolve, 5000));
        
        // 停止加速
        pf.stop();
      }
    }
    
    // 获取公告
    const announcement = pf.getAnnouncement();
    if (announcement) {
      console.log('公告:', announcement);
    }
    
  } catch (error) {
    console.error('错误:', error);
  } finally {
    // 清理资源
    pf.cleanup();
  }
}

// 导出类
module.exports = PlayFast;

// 如果直接运行此文件，执行示例
if (require.main === module) {
  example();
}

