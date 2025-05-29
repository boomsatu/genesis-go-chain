
// Lumina Wallet In-Page Provider
class LuminaWalletProvider {
  constructor() {
    this.isLumina = true;
    this.isMetaMask = true; // For compatibility
    this.chainId = null;
    this.networkVersion = null;
    this.selectedAddress = null;
    this.isConnected = false;
    
    this.requestId = 0;
    this.callbacks = new Map();
    this.eventListeners = new Map();
    
    this.init();
  }

  init() {
    // Listen for responses from content script
    window.addEventListener('message', (event) => {
      if (event.source !== window) return;
      
      if (event.data.type === 'LUMINA_WALLET_RESPONSE') {
        this.handleResponse(event.data.payload);
      }
      
      if (event.data.type === 'LUMINA_WALLET_EVENT') {
        this.handleEvent(event.data.payload);
      }
    });

    // Initialize connection state
    this.request({ method: 'eth_chainId' })
      .then(chainId => {
        this.chainId = chainId;
        this.networkVersion = parseInt(chainId, 16).toString();
      })
      .catch(() => {});

    this.request({ method: 'eth_accounts' })
      .then(accounts => {
        if (accounts.length > 0) {
          this.selectedAddress = accounts[0];
          this.isConnected = true;
        }
      })
      .catch(() => {});
  }

  // EIP-1193 Provider Interface
  async request(args) {
    return new Promise((resolve, reject) => {
      const id = ++this.requestId;
      this.callbacks.set(id, { resolve, reject });
      
      window.postMessage({
        type: 'LUMINA_WALLET_REQUEST',
        payload: {
          id,
          method: args.method,
          params: args.params || []
        }
      }, '*');
    });
  }

  // Legacy methods for compatibility
  async enable() {
    const accounts = await this.request({ method: 'eth_requestAccounts' });
    return accounts;
  }

  async send(method, params = []) {
    return this.request({ method, params });
  }

  async sendAsync(payload, callback) {
    try {
      const result = await this.request({
        method: payload.method,
        params: payload.params
      });
      callback(null, {
        id: payload.id,
        jsonrpc: '2.0',
        result
      });
    } catch (error) {
      callback(error, {
        id: payload.id,
        jsonrpc: '2.0',
        error: {
          code: -32603,
          message: error.message
        }
      });
    }
  }

  // Event handling
  on(event, callback) {
    if (!this.eventListeners.has(event)) {
      this.eventListeners.set(event, []);
    }
    this.eventListeners.get(event).push(callback);
  }

  removeListener(event, callback) {
    const listeners = this.eventListeners.get(event);
    if (listeners) {
      const index = listeners.indexOf(callback);
      if (index > -1) {
        listeners.splice(index, 1);
      }
    }
  }

  emit(event, ...args) {
    const listeners = this.eventListeners.get(event);
    if (listeners) {
      listeners.forEach(callback => callback(...args));
    }
  }

  // Internal methods
  handleResponse(payload) {
    const { id, result, error } = payload;
    const callback = this.callbacks.get(id);
    
    if (callback) {
      this.callbacks.delete(id);
      
      if (error) {
        callback.reject(new Error(error));
      } else {
        callback.resolve(result);
      }
    }
  }

  handleEvent(payload) {
    const { method, params } = payload;
    
    switch (method) {
      case 'chainChanged':
        this.chainId = params[0];
        this.networkVersion = parseInt(params[0], 16).toString();
        this.emit('chainChanged', params[0]);
        this.emit('networkChanged', this.networkVersion);
        break;
        
      case 'accountsChanged':
        const accounts = params[0];
        this.selectedAddress = accounts.length > 0 ? accounts[0] : null;
        this.isConnected = accounts.length > 0;
        this.emit('accountsChanged', accounts);
        break;
        
      case 'connect':
        this.isConnected = true;
        this.emit('connect', { chainId: this.chainId });
        break;
        
      case 'disconnect':
        this.isConnected = false;
        this.selectedAddress = null;
        this.emit('disconnect');
        break;
    }
  }

  // Utility methods
  isAddress(address) {
    return /^0x[a-fA-F0-9]{40}$/.test(address);
  }

  toHex(value) {
    return '0x' + parseInt(value).toString(16);
  }

  fromHex(hex) {
    return parseInt(hex, 16);
  }
}

// Inject the provider into window
window.ethereum = new LuminaWalletProvider();
window.lumina = window.ethereum;

// Announce provider
window.dispatchEvent(new Event('ethereum#initialized'));

// For compatibility with some DApps
Object.defineProperty(window, 'web3', {
  value: {
    currentProvider: window.ethereum
  },
  writable: false
});

// Announce the provider to DApps
const announceProvider = () => {
  const info = {
    uuid: 'lumina-wallet',
    name: 'Lumina Wallet',
    icon: 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIiIGhlaWdodD0iMzIiIHZpZXdCb3g9IjAgMCAzMiAzMiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjMyIiBoZWlnaHQ9IjMyIiByeD0iMTYiIGZpbGw9InVybCgjcGFpbnQwX2xpbmVhcl8xXzEpIi8+CjxwYXRoIGQ9Ik0xNiA4TDI0IDE2TDE2IDI0TDggMTZMMTYgOFoiIGZpbGw9IndoaXRlIi8+CjxkZWZzPgo8bGluZWFyR3JhZGllbnQgaWQ9InBhaW50MF9saW5lYXJfMV8xIiB4MT0iMCIgeTE9IjAiIHgyPSIzMiIgeTI9IjMyIiBncmFkaWVudFVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+CjxzdG9wIHN0b3AtY29sb3I9IiM2NjdFRUEiLz4KPHN0b3Agb2Zmc2V0PSIxIiBzdG9wLWNvbG9yPSIjNzY0QkEyIi8+CjwvbGluZWFyR3JhZGllbnQ+CjwvZGVmcz4KPHN2Zz4K',
    rdns: 'com.lumina.wallet'
  };

  window.dispatchEvent(new CustomEvent('eip6963:announceProvider', {
    detail: Object.freeze({ info, provider: window.ethereum })
  }));
};

announceProvider();

// Listen for provider requests
window.addEventListener('eip6963:requestProvider', announceProvider);
