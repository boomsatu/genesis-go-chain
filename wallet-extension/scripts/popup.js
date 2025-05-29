
class LuminaWallet {
  constructor() {
    this.currentScreen = 'loading';
    this.wallet = null;
    this.init();
  }

  async init() {
    // Check if wallet exists
    const walletData = await this.getStoredWallet();
    
    setTimeout(() => {
      if (walletData) {
        this.showScreen('unlock');
      } else {
        this.showScreen('setup');
      }
    }, 2000);

    this.setupEventListeners();
  }

  setupEventListeners() {
    // Setup events
    document.getElementById('create-wallet').addEventListener('click', () => {
      this.createNewWallet();
    });

    document.getElementById('import-wallet').addEventListener('click', () => {
      this.showScreen('import-wallet-screen');
    });

    document.getElementById('back-from-create').addEventListener('click', () => {
      this.showScreen('setup');
    });

    document.getElementById('back-from-import').addEventListener('click', () => {
      this.showScreen('setup');
    });

    document.getElementById('saved-mnemonic').addEventListener('change', (e) => {
      document.getElementById('confirm-mnemonic').disabled = !e.target.checked;
    });

    document.getElementById('confirm-mnemonic').addEventListener('click', () => {
      this.showScreen('password-screen');
    });

    document.getElementById('set-password').addEventListener('click', () => {
      this.setWalletPassword();
    });

    document.getElementById('import-wallet-btn').addEventListener('click', () => {
      this.importWallet();
    });

    document.getElementById('unlock-wallet').addEventListener('click', () => {
      this.unlockWallet();
    });

    // Main wallet events
    document.getElementById('send-btn').addEventListener('click', () => {
      this.showScreen('send-screen');
    });

    document.getElementById('receive-btn').addEventListener('click', () => {
      this.showScreen('receive-screen');
    });

    document.getElementById('settings-btn').addEventListener('click', () => {
      this.showScreen('settings-screen');
    });

    // Back buttons
    document.getElementById('back-from-send').addEventListener('click', () => {
      this.showScreen('wallet');
    });

    document.getElementById('back-from-receive').addEventListener('click', () => {
      this.showScreen('wallet');
    });

    document.getElementById('back-from-settings').addEventListener('click', () => {
      this.showScreen('wallet');
    });

    // Copy address
    document.getElementById('copy-address').addEventListener('click', () => {
      this.copyToClipboard();
    });

    // Tab switching
    document.querySelectorAll('.tab').forEach(tab => {
      tab.addEventListener('click', (e) => {
        this.switchTab(e.target.dataset.tab);
      });
    });

    // Gas options
    document.querySelectorAll('.gas-option').forEach(option => {
      option.addEventListener('click', (e) => {
        document.querySelectorAll('.gas-option').forEach(o => o.classList.remove('active'));
        e.target.classList.add('active');
      });
    });

    // Send transaction
    document.getElementById('send-transaction').addEventListener('click', () => {
      this.sendTransaction();
    });

    // Settings actions
    document.getElementById('reset-wallet').addEventListener('click', () => {
      this.resetWallet();
    });

    // Network change
    document.getElementById('network-select').addEventListener('change', (e) => {
      this.changeNetwork(e.target.value);
    });
  }

  showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(screen => {
      screen.classList.add('hidden');
    });
    document.getElementById(screenId).classList.remove('hidden');
    this.currentScreen = screenId;
  }

  async createNewWallet() {
    try {
      // Generate mnemonic
      const mnemonic = this.generateMnemonic();
      this.tempMnemonic = mnemonic;
      
      // Display mnemonic
      const mnemonicDisplay = document.getElementById('mnemonic-display');
      mnemonicDisplay.innerHTML = '';
      
      mnemonic.split(' ').forEach((word, index) => {
        const wordElement = document.createElement('div');
        wordElement.className = 'mnemonic-word';
        wordElement.textContent = `${index + 1}. ${word}`;
        mnemonicDisplay.appendChild(wordElement);
      });
      
      this.showScreen('create-wallet-screen');
    } catch (error) {
      console.error('Error creating wallet:', error);
      alert('Error creating wallet');
    }
  }

  async setWalletPassword() {
    const password = document.getElementById('new-password').value;
    const confirmPassword = document.getElementById('confirm-new-password').value;
    
    if (password !== confirmPassword) {
      alert('Passwords do not match');
      return;
    }
    
    if (password.length < 8) {
      alert('Password must be at least 8 characters');
      return;
    }
    
    try {
      // Create wallet from mnemonic
      const wallet = await this.createWalletFromMnemonic(this.tempMnemonic);
      
      // Encrypt and store wallet
      await this.storeWallet(wallet, password);
      
      this.wallet = wallet;
      this.loadWalletData();
      this.showScreen('wallet');
    } catch (error) {
      console.error('Error setting password:', error);
      alert('Error creating wallet');
    }
  }

  async importWallet() {
    const mnemonic = document.getElementById('mnemonic-input').value.trim();
    const password = document.getElementById('password-input').value;
    const confirmPassword = document.getElementById('confirm-password-input').value;
    
    if (password !== confirmPassword) {
      alert('Passwords do not match');
      return;
    }
    
    if (!this.validateMnemonic(mnemonic)) {
      alert('Invalid recovery phrase');
      return;
    }
    
    try {
      const wallet = await this.createWalletFromMnemonic(mnemonic);
      await this.storeWallet(wallet, password);
      
      this.wallet = wallet;
      this.loadWalletData();
      this.showScreen('wallet');
    } catch (error) {
      console.error('Error importing wallet:', error);
      alert('Error importing wallet');
    }
  }

  async unlockWallet() {
    const password = document.getElementById('unlock-password').value;
    
    try {
      const wallet = await this.getStoredWallet();
      const decryptedWallet = await this.decryptWallet(wallet, password);
      
      this.wallet = decryptedWallet;
      this.loadWalletData();
      this.showScreen('wallet');
    } catch (error) {
      console.error('Error unlocking wallet:', error);
      alert('Incorrect password');
    }
  }

  async loadWalletData() {
    if (!this.wallet) return;
    
    // Update address display
    document.getElementById('account-address').textContent = 
      this.wallet.address.slice(0, 6) + '...' + this.wallet.address.slice(-4);
    document.getElementById('receive-address').textContent = this.wallet.address;
    
    // Load balance
    await this.updateBalance();
    
    // Load tokens
    await this.loadTokens();
    
    // Load activities
    await this.loadActivities();
    
    // Generate QR code
    this.generateQRCode();
  }

  async updateBalance() {
    try {
      // Get balance from blockchain
      const balance = await this.getBalance(this.wallet.address);
      document.getElementById('balance-amount').textContent = balance.toFixed(4);
      
      // Convert to USD (mock conversion)
      const usdValue = (balance * 100).toFixed(2); // Mock LUM to USD rate
      document.getElementById('balance-usd').textContent = `$${usdValue} USD`;
    } catch (error) {
      console.error('Error updating balance:', error);
    }
  }

  async loadTokens() {
    const tokenList = document.getElementById('token-list');
    tokenList.innerHTML = '';
    
    // Add native token
    const nativeToken = this.createTokenElement('LUM', 'Lumina', '100.00', '$10,000.00');
    tokenList.appendChild(nativeToken);
    
    // Load other tokens (mock data)
    const tokens = [
      { symbol: 'USDC', name: 'USD Coin', balance: '1,000.00', value: '$1,000.00' },
      { symbol: 'WETH', name: 'Wrapped Ether', balance: '0.5', value: '$1,250.00' }
    ];
    
    tokens.forEach(token => {
      const tokenElement = this.createTokenElement(token.symbol, token.name, token.balance, token.value);
      tokenList.appendChild(tokenElement);
    });
  }

  createTokenElement(symbol, name, balance, value) {
    const tokenItem = document.createElement('div');
    tokenItem.className = 'token-item';
    tokenItem.innerHTML = `
      <div class="token-icon">${symbol.charAt(0)}</div>
      <div class="token-info">
        <div class="token-name">${name}</div>
        <div class="token-symbol">${symbol}</div>
      </div>
      <div class="token-balance">
        <div class="token-amount">${balance}</div>
        <div class="token-value">${value}</div>
      </div>
    `;
    return tokenItem;
  }

  async loadActivities() {
    const activityList = document.getElementById('activity-list');
    activityList.innerHTML = '';
    
    // Mock activity data
    const activities = [
      { type: 'Send', icon: '📤', time: '2 hours ago', amount: '-1.5 LUM' },
      { type: 'Receive', icon: '📥', time: '1 day ago', amount: '+10.0 LUM' },
      { type: 'Swap', icon: '🔄', time: '3 days ago', amount: '5.0 LUM → 500 USDC' }
    ];
    
    activities.forEach(activity => {
      const activityItem = document.createElement('div');
      activityItem.className = 'activity-item';
      activityItem.innerHTML = `
        <div class="activity-icon">${activity.icon}</div>
        <div class="activity-info">
          <div class="activity-type">${activity.type}</div>
          <div class="activity-time">${activity.time}</div>
        </div>
        <div class="activity-amount">${activity.amount}</div>
      `;
      activityList.appendChild(activityItem);
    });
  }

  switchTab(tabName) {
    document.querySelectorAll('.tab').forEach(tab => {
      tab.classList.remove('active');
    });
    document.querySelectorAll('.tab-panel').forEach(panel => {
      panel.classList.add('hidden');
    });
    
    document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');
    document.getElementById(`${tabName}-tab`).classList.remove('hidden');
  }

  async sendTransaction() {
    const to = document.getElementById('send-to').value;
    const amount = document.getElementById('send-amount').value;
    const token = document.getElementById('send-token').value;
    
    if (!to || !amount) {
      alert('Please fill in all fields');
      return;
    }
    
    if (!this.isValidAddress(to)) {
      alert('Invalid recipient address');
      return;
    }
    
    try {
      // Send transaction to blockchain
      const txHash = await this.submitTransaction(to, amount, token);
      alert(`Transaction sent! Hash: ${txHash}`);
      this.showScreen('wallet');
      this.updateBalance();
    } catch (error) {
      console.error('Error sending transaction:', error);
      alert('Transaction failed');
    }
  }

  generateQRCode() {
    const qrCode = document.getElementById('qr-code');
    qrCode.innerHTML = `
      <div style="font-size: 12px; color: #6b7280;">
        QR Code for<br>${this.wallet.address}
      </div>
    `;
  }

  copyToClipboard() {
    navigator.clipboard.writeText(this.wallet.address);
    alert('Address copied to clipboard!');
  }

  changeNetwork(network) {
    console.log('Switching to network:', network);
    // Implement network switching logic
    this.updateBalance();
  }

  resetWallet() {
    if (confirm('Are you sure you want to reset your wallet? This action cannot be undone.')) {
      chrome.storage.local.clear();
      this.wallet = null;
      this.showScreen('setup');
    }
  }

  // Utility functions
  generateMnemonic() {
    const words = [
      'abandon', 'ability', 'able', 'about', 'above', 'absent', 'absorb', 'abstract',
      'absurd', 'abuse', 'access', 'accident', 'account', 'accuse', 'achieve', 'acid',
      'acoustic', 'acquire', 'across', 'act', 'action', 'actor', 'actress', 'actual'
    ];
    
    const mnemonic = [];
    for (let i = 0; i < 12; i++) {
      mnemonic.push(words[Math.floor(Math.random() * words.length)]);
    }
    return mnemonic.join(' ');
  }

  validateMnemonic(mnemonic) {
    const words = mnemonic.trim().split(' ');
    return words.length === 12;
  }

  async createWalletFromMnemonic(mnemonic) {
    // Mock wallet creation - in production, use proper crypto libraries
    const privateKey = this.hashString(mnemonic);
    const address = '0x' + this.hashString(privateKey).slice(0, 40);
    
    return {
      mnemonic,
      privateKey,
      address
    };
  }

  hashString(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      const char = str.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash;
    }
    return Math.abs(hash).toString(16).padStart(32, '0');
  }

  async storeWallet(wallet, password) {
    // In production, properly encrypt the wallet
    const encryptedWallet = btoa(JSON.stringify(wallet) + password);
    await chrome.storage.local.set({ wallet: encryptedWallet });
  }

  async getStoredWallet() {
    const result = await chrome.storage.local.get(['wallet']);
    return result.wallet;
  }

  async decryptWallet(encryptedWallet, password) {
    try {
      const decrypted = atob(encryptedWallet);
      const wallet = JSON.parse(decrypted.replace(password, ''));
      return wallet;
    } catch (error) {
      throw new Error('Invalid password');
    }
  }

  async getBalance(address) {
    // Mock balance - in production, call blockchain RPC
    return Math.random() * 1000;
  }

  async submitTransaction(to, amount, token) {
    // Mock transaction - in production, sign and submit to blockchain
    return '0x' + Math.random().toString(16).slice(2, 66);
  }

  isValidAddress(address) {
    return /^0x[a-fA-F0-9]{40}$/.test(address);
  }
}

// Initialize wallet when popup opens
document.addEventListener('DOMContentLoaded', () => {
  new LuminaWallet();
});
