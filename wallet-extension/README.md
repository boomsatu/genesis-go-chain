
# 🌟 Lumina Wallet - Chrome Extension

A professional cryptocurrency wallet extension for Lumina Blockchain with MetaMask-compatible features.

## ✨ Features

### 🔐 Security
- **HD Wallet**: Hierarchical Deterministic wallet with BIP39 mnemonic phrases
- **Password Protection**: AES-256 encryption for wallet storage
- **Auto-lock**: Automatic wallet locking after inactivity
- **Secure Random**: Cryptographically secure random number generation

### 💰 Wallet Management
- **Multi-Account**: Support for multiple accounts
- **Token Support**: Native LUM and custom ERC-20 tokens
- **NFT Support**: View and manage NFT collections
- **Balance Tracking**: Real-time balance updates and USD conversion

### 🌐 Network Support
- **Multi-Network**: Support for mainnet, testnet, and custom networks
- **RPC Configuration**: Custom RPC endpoint configuration
- **Network Switching**: Easy network switching with visual indicators
- **Gas Optimization**: Smart gas fee estimation and optimization

### 🔗 DApp Integration
- **Web3 Provider**: Full EIP-1193 provider interface
- **MetaMask Compatibility**: Compatible with existing DApps
- **EIP-6963**: Multi-wallet discovery standard support
- **Custom Events**: Real-time event notifications to DApps

### 📱 User Interface
- **Modern Design**: Clean and intuitive user interface
- **Responsive Layout**: Optimized for different screen sizes
- **Dark/Light Theme**: Automatic theme detection
- **Accessibility**: WCAG 2.1 compliance

## 🚀 Installation

### For Users

1. **Download**: Download the latest release from GitHub
2. **Extract**: Extract the zip file to a folder
3. **Chrome Extensions**: Go to `chrome://extensions/`
4. **Developer Mode**: Enable "Developer mode" (top right)
5. **Load Extension**: Click "Load unpacked" and select the extracted folder
6. **Pin Extension**: Pin the extension to your toolbar for easy access

### For Developers

```bash
# Clone the repository
git clone <repository-url>
cd blockchain-node/wallet-extension

# Install dependencies
npm install

# Build icons and assets
npm run build:icons

# Load as unpacked extension in Chrome
npm run dev
```

## 🛠 Development

### Project Structure

```
wallet-extension/
├── manifest.json          # Extension manifest
├── popup.html             # Main wallet interface
├── background.js          # Service worker
├── content.js             # Content script
├── inpage.js              # In-page provider
├── styles/
│   └── popup.css          # Styles for popup
├── scripts/
│   └── popup.js           # Popup functionality
├── icons/                 # Extension icons
└── README.md              # Documentation
```

### Key Components

#### 1. Background Script (`background.js`)
- **Service Worker**: Handles extension lifecycle and persistence
- **Message Handling**: Processes messages from popup and content scripts
- **Wallet Management**: Manages wallet state and cryptographic operations
- **Network Communication**: Handles blockchain RPC calls

#### 2. Content Script (`content.js`)
- **DApp Bridge**: Bridges communication between DApps and extension
- **Event Handling**: Forwards events between in-page script and background
- **Security**: Validates and sanitizes messages

#### 3. In-Page Provider (`inpage.js`)
- **Web3 Provider**: Implements EIP-1193 provider interface
- **MetaMask Compatibility**: Provides MetaMask-compatible API
- **Event Emission**: Emits wallet events to DApps

#### 4. Popup Interface (`popup.html` + `popup.js`)
- **User Interface**: Main wallet interface for users
- **State Management**: Manages UI state and user interactions
- **Wallet Operations**: Handles wallet creation, import, and management

### 🔧 Configuration

#### Network Configuration
```javascript
const networks = {
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
};
```

#### Security Settings
```javascript
const settings = {
  autoLock: 15,          // Auto-lock timeout in minutes
  currency: 'USD',       // Display currency
  language: 'en',        // Interface language
  encryption: 'AES-256'  // Encryption algorithm
};
```

## 🔐 Security Features

### Wallet Security
- **BIP39 Mnemonics**: Standard 12-word recovery phrases
- **HD Derivation**: BIP44 hierarchical deterministic key derivation
- **Password Encryption**: AES-256 encryption for stored wallets
- **Auto-lock**: Configurable auto-lock timeout

### DApp Security
- **Origin Validation**: Validates requesting DApp origins
- **Permission System**: Granular permission management
- **User Confirmation**: User approval required for transactions
- **Secure Context**: Requires HTTPS for sensitive operations

### Data Protection
- **Local Storage**: Encrypted local storage only
- **No Remote Logging**: No data sent to external servers
- **Secure Random**: Cryptographically secure randomness
- **Memory Protection**: Sensitive data cleared from memory

## 🌐 DApp Integration

### Basic Integration

```javascript
// Request wallet connection
const accounts = await window.ethereum.request({
  method: 'eth_requestAccounts'
});

// Get current network
const chainId = await window.ethereum.request({
  method: 'eth_chainId'
});

// Send transaction
const txHash = await window.ethereum.request({
  method: 'eth_sendTransaction',
  params: [{
    from: accounts[0],
    to: '0x...',
    value: '0x1bc16d674ec80000', // 2 ETH in wei
    gas: '0x5208',              // 21000 gas
    gasPrice: '0x9184e72a000'   // 10 Gwei
  }]
});
```

### Event Handling

```javascript
// Listen for account changes
window.ethereum.on('accountsChanged', (accounts) => {
  console.log('Accounts changed:', accounts);
});

// Listen for network changes
window.ethereum.on('chainChanged', (chainId) => {
  console.log('Network changed:', chainId);
});

// Listen for connection events
window.ethereum.on('connect', (connectInfo) => {
  console.log('Connected:', connectInfo);
});
```

### Advanced Features

```javascript
// Sign typed data
const signature = await window.ethereum.request({
  method: 'eth_signTypedData_v4',
  params: [account, typedData]
});

// Add custom network
await window.ethereum.request({
  method: 'wallet_addEthereumChain',
  params: [{
    chainId: '0x539',
    chainName: 'Lumina Testnet',
    rpcUrls: ['http://localhost:8546'],
    nativeCurrency: {
      name: 'Lumina',
      symbol: 'LUM',
      decimals: 18
    }
  }]
});
```

## 🧪 Testing

### Manual Testing

1. **Wallet Creation**: Test wallet creation with mnemonic generation
2. **Import Wallet**: Test wallet import with existing mnemonic
3. **Transaction Sending**: Test sending transactions with gas estimation
4. **DApp Connection**: Test connection to sample DApp
5. **Network Switching**: Test switching between networks

### Automated Testing

```bash
# Run tests
npm test

# Run linting
npm run lint

# Build for production
npm run build
```

### Test DApp

Create a simple test page to verify wallet functionality:

```html
<!DOCTYPE html>
<html>
<head>
  <title>Lumina Wallet Test</title>
</head>
<body>
  <button id="connect">Connect Wallet</button>
  <button id="send">Send Transaction</button>
  
  <script>
    document.getElementById('connect').onclick = async () => {
      const accounts = await window.ethereum.request({
        method: 'eth_requestAccounts'
      });
      console.log('Connected:', accounts);
    };
    
    document.getElementById('send').onclick = async () => {
      // Test transaction
    };
  </script>
</body>
</html>
```

## 📦 Building for Production

### Build Process

```bash
# Install dependencies
npm install

# Generate icons
npm run build:icons

# Create distribution package
npm run build:zip
```

### Distribution

1. **Zip Package**: Creates `lumina-wallet.zip` for Chrome Web Store
2. **Icon Generation**: Generates all required icon sizes
3. **Manifest Validation**: Validates manifest.json format
4. **Code Minification**: Minifies JavaScript for production

### Chrome Web Store Submission

1. **Developer Account**: Register Chrome Web Store developer account
2. **App Listing**: Create app listing with screenshots and description
3. **Privacy Policy**: Provide privacy policy and data usage disclosure
4. **Review Process**: Submit for Google review (typically 2-3 days)

## 🔧 Troubleshooting

### Common Issues

#### Extension Not Loading
- **Check Manifest**: Verify manifest.json is valid
- **File Permissions**: Ensure all files are readable
- **Developer Mode**: Confirm developer mode is enabled

#### DApp Connection Issues
- **Content Script**: Verify content script injection
- **HTTPS Context**: Some features require HTTPS
- **Cross-Origin**: Check for CORS issues

#### Transaction Failures
- **Gas Estimation**: Verify gas estimation is working
- **Network Connection**: Check RPC endpoint connectivity
- **Account Balance**: Ensure sufficient balance for transactions

### Debug Mode

Enable debug logging by setting:

```javascript
localStorage.setItem('lumina-debug', 'true');
```

### Support

For technical support:
- **GitHub Issues**: Report bugs and feature requests
- **Documentation**: Check comprehensive documentation
- **Community**: Join Discord community for help

## 🚀 Future Enhancements

### Planned Features

- **Hardware Wallet**: Ledger and Trezor integration
- **Multi-Signature**: Multi-sig wallet support
- **Staking**: Native staking interface
- **DeFi Integration**: Built-in DeFi protocol support
- **Mobile App**: React Native mobile application

### Roadmap

- **v1.1**: Hardware wallet support
- **v1.2**: Advanced DeFi features
- **v1.3**: Mobile companion app
- **v2.0**: Multi-chain support

---

**Built with ❤️ for the Lumina Blockchain ecosystem**

For more information, visit our [documentation](../docs/) or join our [community](https://discord.gg/lumina).
