import { getIcon } from '@iconify/vue';

function isUrl(value) {
  return !!value && (value.startsWith('http://') || value.startsWith('https://'));
}

function isIconifyIcon(value) {
  return !!value && value.includes(':');
}

function setFavicon(iconUrl) {
  const link = document.querySelector("link[rel~='icon']");
  if (link) {
    link.href = iconUrl;
    return;
  }

  const newLink = document.createElement('link');
  newLink.rel = 'icon';
  newLink.href = iconUrl;
  document.head.appendChild(newLink);
}

function setIconifyIcon(iconName) {
  try {
    const iconData = getIcon(iconName);
    if (!iconData) return;

    const svg = `
      <svg xmlns="http://www.w3.org/2000/svg"
           width="64" height="64"
           viewBox="0 0 ${iconData.width} ${iconData.height}"
           preserveAspectRatio="xMidYMid meet">
        <g transform="translate(0, 4)">
          ${iconData.body}
        </g>
      </svg>
    `;

    setFavicon(`data:image/svg+xml,${encodeURIComponent(svg)}`);
  } catch (error) {
    console.error('Error loading icon:', error);
  }
}

export function applySiteMeta(settings) {
  const icon = settings && settings.icon;

  if (isUrl(icon)) {
    setFavicon(icon);
  } else if (isIconifyIcon(icon)) {
    setIconifyIcon(icon);
  } else {
    setFavicon('https://img1.baidu.com/it/u=1217061905,2277984247&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500');
  }

  document.title = settings && settings.title ? settings.title : 'ZNav';
}
