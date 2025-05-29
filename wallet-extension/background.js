
// Lumina Wallet Background Script
class LuminaWalletBackground {
  constructor() {
    this.init();
  }

  init() {
    // Handle extension installation
    chrome.runtime.onInstalled.addListener((details) => {
      if (details.reason === 'install') {
        console.log('Lumina Wallet installed');
        this.setupInitialState();
      }
    });

    // Handle messages from content scripts and popup
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
      this.handleMessage(request, sender, sendResponse);
      return true; // Keep message channel open for async responses
    });

    // Handle web3 provider requests
    chrome.runtime.onMessageExternal.addListener((request, sender, sendResponse) => {
      this.handleExternalMessage(request, sender, sendResponse);
      return true;
    });

    // Monitor network changes
    this.setupNetworkMonitoring();
  }

  setupInitialState() {
    // Set default configuration
    chrome.storage.local.set({
      networks: {
        'lumina-mainnet': {
          name: 'Lumina Mainnet',
          rpcUrl: 'http://localhost:8545',
          chainId: 1337,
          symbol: 'LUM',
          blockExplorer: 'http://localhost:3000'
        },
        'lumina-testnet': {
          name: 'Lumina Testnet',
          rpcUrl: 'http://localhost:8546',
          chainId: 1338,
          symbol: 'LUM',
          blockExplorer: 'http://localhost:3001'
        }
      },
      currentNetwork: 'lumina-mainnet',
      connectedSites: {},
      settings: {
        autoLock: 15, // minutes
        currency: 'USD',
        language: 'en'
      }
    });
  }

  async handleMessage(request, sender, sendResponse) {
    const { type, data } = request;

    try {
      switch (type) {
        case 'GET_WALLET_STATE':
          const walletState = await this.getWalletState();
          sendResponse({ success: true, data: walletState });
          break;

        case 'CONNECT_WALLET':
          const connection = await this.connectWallet(data.origin);
          sendResponse({ success: true, data: connection });
          break;

        case 'DISCONNECT_WALLET':
          await this.disconnectWallet(data.origin);
          sendResponse({ success: true });
          break;

        case 'SIGN_TRANSACTION':
          const signedTx = await this.signTransaction(data);
          sendResponse({ success: true, data: signedTx });
          break;

        case 'SEND_TRANSACTION':
          const txHash = await this.sendTransaction(data);
          sendResponse({ success: true, data: txHash });
          break;

        case 'GET_BALANCE':
          const balance = await this.getBalance(data.address, data.network);
          sendResponse({ success: true, data: balance });
          break;

        case 'GET_TRANSACTION_HISTORY':
          const history = await this.getTransactionHistory(data.address, data.network);
          sendResponse({ success: true, data: history });
          break;

        case 'ADD_NETWORK':
          await this.addNetwork(data);
          sendResponse({ success: true });
          break;

        case 'SWITCH_NETWORK':
          await this.switchNetwork(data.networkId);
          sendResponse({ success: true });
          break;

        case 'ADD_TOKEN':
          await this.addToken(data);
          sendResponse({ success: true });
          break;

        default:
          sendResponse({ success: false, error: 'Unknown message type' });
      }
    } catch (error) {
      console.error('Error handling message:', error);
      sendResponse({ success: false, error: error.message });
    }
  }

  async handleExternalMessage(request, sender, sendResponse) {
    // Handle requests from DApps
    const { method, params } = request;
    const origin = sender.origin || sender.url;

    try {
      // Check if site is connected
      const isConnected = await this.isSiteConnected(origin);
      
      if (!isConnected && method !== 'eth_requestAccounts') {
        sendResponse({ error: 'Unauthorized' });
        return;
      }

      switch (method) {
        case 'eth_requestAccounts':
          const accounts = await this.requestAccounts(origin);
          sendResponse({ result: accounts });
          break;

        case 'eth_accounts':
          const connectedAccounts = await this.getConnectedAccounts(origin);
          sendResponse({ result: connectedAccounts });
          break;

        case 'eth_chainId':
          const chainId = await this.getChainId();
          sendResponse({ result: chainId });
          break;

        case 'eth_getBalance':
          const balance = await this.getBalance(params[0], params[1]);
          sendResponse({ result: balance });
          break;

        case 'eth_sendTransaction':
          const txHash = await this.sendTransaction(params[0]);
          sendResponse({ result: txHash });
          break;

        case 'eth_signTransaction':
          const signedTx = await this.signTransaction(params[0]);
          sendResponse({ result: signedTx });
          break;

        case 'personal_sign':
          const signature = await this.personalSign(params[0], params[1]);
          sendResponse({ result: signature });
          break;

        case 'eth_signTypedData_v4':
          const typedSignature = await this.signTypedData(params[0], params[1]);
          sendResponse({ result: typedSignature });
          break;

        case 'wallet_addEthereumChain':
          await this.addEthereumChain(params[0]);
          sendResponse({ result: null });
          break;

        case 'wallet_switchEthereumChain':
          await this.switchEthereumChain(params[0]);
          sendResponse({ result: null });
          break;

        default:
          sendResponse({ error: 'Method not supported' });
      }
    } catch (error) {
      console.error('Error handling external message:', error);
      sendResponse({ error: error.message });
    }
  }

  setupNetworkMonitoring() {
    // Monitor network changes and notify connected sites
    chrome.storage.onChanged.addListener((changes, namespace) => {
      if (namespace === 'local' && changes.currentNetwork) {
        this.notifyNetworkChange(changes.currentNetwork.newValue);
      }
    });
  }

  async getWalletState() {
    const result = await chrome.storage.local.get(['wallet', 'currentNetwork', 'networks']);
    return {
      isUnlocked: !!result.wallet,
      currentNetwork: result.currentNetwork,
      networks: result.networks
    };
  }

  async connectWallet(origin) {
    // Check if already connected
    const connectedSites = await this.getConnectedSites();
    if (connectedSites[origin]) {
      return { connected: true, accounts: connectedSites[origin].accounts };
    }

    // Request user approval
    const approval = await this.requestUserApproval(origin, 'connect');
    if (!approval) {
      throw new Error('User rejected connection');
    }

    // Get current account
    const wallet = await this.getCurrentWallet();
    const accounts = [wallet.address];

    // Store connection
    connectedSites[origin] = {
      accounts,
      permissions: ['eth_accounts'],
      connected: Date.now()
    };
    await chrome.storage.local.set({ connectedSites });

    return { connected: true, accounts };
  }

  async disconnectWallet(origin) {
    const connectedSites = await this.getConnectedSites();
    delete connectedSites[origin];
    await chrome.storage.local.set({ connectedSites });
  }

  async signTransaction(txData) {
    // Request user approval
    const approval = await this.requestUserApproval(null, 'signTransaction', txData);
    if (!approval) {
      throw new Error('User rejected transaction');
    }

    // Sign transaction with wallet
    const wallet = await this.getCurrentWallet();
    return this.cryptoSign(txData, wallet.privateKey);
  }

  async sendTransaction(txData) {
    // Sign transaction
    const signedTx = await this.signTransaction(txData);
    
    // Send to network
    const network = await this.getCurrentNetwork();
    const response = await fetch(network.rpcUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        jsonrpc: '2.0',
        method: 'eth_sendRawTransaction',
        params: [signedTx],
        id: 1
      })
    });

    const result = await response.json();
    if (result.error) {
      throw new Error(result.error.message);
    }

    return result.result;
  }

  async getBalance(address, blockTag = 'latest') {
    const network = await this.getCurrentNetwork();
    const response = await fetch(network.rpcUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        jsonrpc: '2.0',
        method: 'eth_getBalance',
        params: [address, blockTag],
        id: 1
      })
    });

    const result = await response.json();
    if (result.error) {
      throw new Error(result.error.message);
    }

    return result.result;
  }

  async getTransactionHistory(address, network) {
    // Mock implementation - in production, fetch from blockchain
    return [
      {
        hash: '0x1234...',
        from: address,
        to: '0x5678...',
        value: '1000000000000000000',
        timestamp: Date.now() - 3600000,
        status: 'success'
      }
    ];
  }

  async addNetwork(networkData) {
    const networks = await this.getNetworks();
    networks[networkData.chainId] = networkData;
    await chrome.storage.local.set({ networks });
  }

  async switchNetwork(networkId) {
    const networks = await this.getNetworks();
    if (!networks[networkId]) {
      throw new Error('Network not found');
    }
    await chrome.storage.local.set({ currentNetwork: networkId });
  }

  async addToken(tokenData) {
    const tokens = await this.getTokens();
    tokens[tokenData.address] = tokenData;
    await chrome.storage.local.set({ tokens });
  }

  async requestAccounts(origin) {
    const approval = await this.requestUserApproval(origin, 'requestAccounts');
    if (!approval) {
      throw new Error('User rejected request');
    }

    const connection = await this.connectWallet(origin);
    return connection.accounts;
  }

  async getConnectedAccounts(origin) {
    const connectedSites = await this.getConnectedSites();
    if (!connectedSites[origin]) {
      return [];
    }
    return connectedSites[origin].accounts;
  }

  async getChainId() {
    const network = await this.getCurrentNetwork();
    return '0x' + network.chainId.toString(16);
  }

  async personalSign(message, address) {
    const approval = await this.requestUserApproval(null, 'personalSign', { message, address });
    if (!approval) {
      throw new Error('User rejected signing');
    }

    const wallet = await this.getCurrentWallet();
    return this.cryptoSign(message, wallet.privateKey);
  }

  async signTypedData(address, typedData) {
    const approval = await this.requestUserApproval(null, 'signTypedData', { address, typedData });
    if (!approval) {
      throw new Error('User rejected signing');
    }

    const wallet = await this.getCurrentWallet();
    return this.cryptoSign(typedData, wallet.privateKey);
  }

  async addEthereumChain(chainData) {
    const approval = await this.requestUserApproval(null, 'addChain', chainData);
    if (!approval) {
      throw new Error('User rejected chain addition');
    }

    await this.addNetwork({
      chainId: parseInt(chainData.chainId, 16),
      name: chainData.chainName,
      rpcUrl: chainData.rpcUrls[0],
      symbol: chainData.nativeCurrency.symbol,
      blockExplorer: chainData.blockExplorerUrls?.[0]
    });
  }

  async switchEthereumChain(chainData) {
    const chainId = parseInt(chainData.chainId, 16);
    await this.switchNetwork(chainId);
  }

  // Utility methods
  async getCurrentWallet() {
    const result = await chrome.storage.local.get(['wallet']);
    if (!result.wallet) {
      throw new Error('Wallet not found');
    }
    // In production, decrypt wallet here
    return JSON.parse(atob(result.wallet));
  }

  async getCurrentNetwork() {
    const result = await chrome.storage.local.get(['currentNetwork', 'networks']);
    return result.networks[result.currentNetwork];
  }

  async getConnectedSites() {
    const result = await chrome.storage.local.get(['connectedSites']);
    return result.connectedSites || {};
  }

  async getNetworks() {
    const result = await chrome.storage.local.get(['networks']);
    return result.networks || {};
  }

  async getTokens() {
    const result = await chrome.storage.local.get(['tokens']);
    return result.tokens || {};
  }

  async isSiteConnected(origin) {
    const connectedSites = await this.getConnectedSites();
    return !!connectedSites[origin];
  }

  async requestUserApproval(origin, action, data = null) {
    // In production, show approval popup to user
    // For now, return true (auto-approve)
    return true;
  }

  async notifyNetworkChange(newNetwork) {
    // Notify all connected sites about network change
    const connectedSites = await this.getConnectedSites();
    Object.keys(connectedSites).forEach(origin => {
      // Send message to content script to notify the site
      chrome.tabs.query({ url: `${origin}/*` }, (tabs) => {
        tabs.forEach(tab => {
          chrome.tabs.sendMessage(tab.id, {
            type: 'NETWORK_CHANGED',
            data: { networkId: newNetwork }
          });
        });
      });
    });
  }

  cryptoSign(data, privateKey) {
    // Mock signing - in production, use proper crypto libraries
    return '0x' + Math.random().toString(16).slice(2, 130);
  }
}

// Initialize background script
new LuminaWalletBackground();
