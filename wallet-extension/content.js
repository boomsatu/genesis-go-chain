
// Lumina Wallet Content Script
class LuminaWalletContent {
  constructor() {
    this.init();
  }

  init() {
    // Inject the in-page script
    this.injectInPageScript();
    
    // Listen for messages from the in-page script
    window.addEventListener('message', (event) => {
      if (event.source !== window || !event.data.type) return;
      
      if (event.data.type === 'LUMINA_WALLET_REQUEST') {
        this.handleWalletRequest(event.data);
      }
    });

    // Listen for messages from background script
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
      if (request.type === 'NETWORK_CHANGED') {
        this.notifyNetworkChange(request.data);
      }
      if (request.type === 'ACCOUNTS_CHANGED') {
        this.notifyAccountsChange(request.data);
      }
    });
  }

  injectInPageScript() {
    const script = document.createElement('script');
    script.src = chrome.runtime.getURL('inpage.js');
    script.onload = () => script.remove();
    (document.head || document.documentElement).appendChild(script);
  }

  async handleWalletRequest(data) {
    const { id, method, params } = data.payload;

    try {
      // Forward request to background script
      const response = await chrome.runtime.sendMessage({
        type: 'WALLET_REQUEST',
        data: { method, params, origin: window.location.origin }
      });

      // Send response back to in-page script
      window.postMessage({
        type: 'LUMINA_WALLET_RESPONSE',
        payload: {
          id,
          result: response.success ? response.data : null,
          error: response.success ? null : response.error
        }
      }, '*');
    } catch (error) {
      window.postMessage({
        type: 'LUMINA_WALLET_RESPONSE',
        payload: {
          id,
          result: null,
          error: error.message
        }
      }, '*');
    }
  }

  notifyNetworkChange(data) {
    window.postMessage({
      type: 'LUMINA_WALLET_EVENT',
      payload: {
        method: 'chainChanged',
        params: ['0x' + data.networkId.toString(16)]
      }
    }, '*');
  }

  notifyAccountsChange(data) {
    window.postMessage({
      type: 'LUMINA_WALLET_EVENT',
      payload: {
        method: 'accountsChanged',
        params: [data.accounts]
      }
    }, '*');
  }
}

// Initialize content script
new LuminaWalletContent();
