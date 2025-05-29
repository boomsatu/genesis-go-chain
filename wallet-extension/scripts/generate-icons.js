
// Generate icons for Lumina Wallet Extension
const fs = require('fs');
const path = require('path');

// Create icons directory if it doesn't exist
const iconsDir = path.join(__dirname, '..', 'icons');
if (!fs.existsSync(iconsDir)) {
  fs.mkdirSync(iconsDir, { recursive: true });
}

// SVG icon template
const createSvgIcon = (size) => `
<svg width="${size}" height="${size}" viewBox="0 0 ${size} ${size}" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="${size}" height="${size}" rx="${size/2}" fill="url(#paint0_linear_1_1)"/>
  <path d="M${size/2} ${size/4}L${size*3/4} ${size/2}L${size/2} ${size*3/4}L${size/4} ${size/2}L${size/2} ${size/4}Z" fill="white"/>
  <defs>
    <linearGradient id="paint0_linear_1_1" x1="0" y1="0" x2="${size}" y2="${size}" gradientUnits="userSpaceOnUse">
      <stop stop-color="#667EEA"/>
      <stop offset="1" stop-color="#764BA2"/>
    </linearGradient>
  </defs>
</svg>
`;

// Generate PNG placeholder files (in production, use proper image generation)
const createPngPlaceholder = (size) => `
data:image/svg+xml;base64,${Buffer.from(createSvgIcon(size)).toString('base64')}
`;

const iconSizes = [16, 32, 48, 128];

iconSizes.forEach(size => {
  const svgContent = createSvgIcon(size);
  fs.writeFileSync(path.join(iconsDir, `icon${size}.svg`), svgContent.trim());
  
  // Create placeholder PNG info file
  const pngInfo = {
    size: size,
    format: 'PNG',
    description: `Lumina Wallet icon ${size}x${size}`,
    placeholder: createPngPlaceholder(size)
  };
  
  fs.writeFileSync(
    path.join(iconsDir, `icon${size}.json`), 
    JSON.stringify(pngInfo, null, 2)
  );
});

console.log('Icons generated successfully!');
console.log('Note: In production, convert SVG files to PNG using a proper image processing tool.');
