function initUserquakeTimeline(objectId) {
  const userquakeContainer = document.querySelector(`[data-userquake-id="${objectId}"]`);
  if (!userquakeContainer) return;

  const playBtn = userquakeContainer.querySelector('.timeline-play-btn');
  const slider = userquakeContainer.querySelector('.timeline-slider');
  const speedSelect = userquakeContainer.querySelector('.timeline-speed-select');
  const currentTimeDisplay = userquakeContainer.querySelector('.timeline-current-time');
  const confidenceDisplay = userquakeContainer.querySelector('.timeline-confidence-display');
  const timelineImage = userquakeContainer.querySelector('.timeline-image');
  const timelineImageLink = userquakeContainer.querySelector('.timeline-image-link');

  let timeseriesData = [];
  let isPlaying = false;
  let currentIndex = 0;
  let animationInterval = null;
  let startTime = null;
  let totalDurationSeconds = 0;
  let currentTimeSeconds = 0;
  let preloadedImages = new Map();
  let isPreloading = false;
  let hasPreloaded = false;

  fetch(`/api/timeseries/${objectId}`)
    .then(response => response.json())
    .then(data => {
      timeseriesData = data;
      if (data.length > 0) {
        let confidenceStartTime = null;
        for (let i = 0; i < data.length; i++) {
          if (data[i].confidence && data[i].confidence > 0.9) {
            confidenceStartTime = new Date(data[i].updated_at.replace(/\//g, '-').replace(' ', 'T'));
            break;
          }
        }
        
        startTime = confidenceStartTime || new Date(data[0].started_at.replace(/\//g, '-').replace(' ', 'T'));
        const endTime = new Date(data[data.length - 1].updated_at.replace(/\//g, '-').replace(' ', 'T'));
        totalDurationSeconds = Math.ceil((endTime - startTime) / 1000);
        
        slider.min = 0;
        slider.max = totalDurationSeconds;
        slider.value = totalDurationSeconds;
        
        currentTimeSeconds = totalDurationSeconds;
        currentIndex = data.length - 1;
        updateDisplay();
      }
    })
    .catch(error => console.error('Error loading timeseries data:', error));

  function updateDisplay() {
    if (timeseriesData.length === 0) return;
    
    currentIndex = findDataIndexForTime(currentTimeSeconds);
    const current = timeseriesData[currentIndex];
    if (!current) return;

    const updatedAt = new Date(current.updated_at.replace(/\//g, '-').replace(' ', 'T'));
    currentTimeDisplay.textContent = updatedAt.toLocaleTimeString('ja-JP');

    const processed = processUserquakeData(current);
    updateConfidenceDisplay(processed);
    
    updateImage(current);
  }

  function findDataIndexForTime(targetSeconds) {
    if (!startTime || timeseriesData.length === 0) return 0;
    
    const targetTime = new Date(startTime.getTime() + targetSeconds * 1000);
    
    let bestIndex = 0;
    let minDiff = Infinity;
    
    for (let i = 0; i < timeseriesData.length; i++) {
      const dataTime = new Date(timeseriesData[i].updated_at.replace(/\//g, '-').replace(' ', 'T'));
      const diff = Math.abs(dataTime - targetTime);
      
      if (diff < minDiff) {
        minDiff = diff;
        bestIndex = i;
      }
      
      if (dataTime > targetTime) {
        break;
      }
    }
    
    return bestIndex;
  }


  function processUserquakeData(data) {
    const areaConfidences = data.area_confidences || {};
    
    let max = 0.125;
    for (const [area, conf] of Object.entries(areaConfidences)) {
      if (conf.confidence > max) max = conf.confidence;
    }
    
    const factor = 1.0 / max;
    const normalizedAreas = [];
    
    for (const [area, conf] of Object.entries(areaConfidences)) {
      const normalizedConf = conf.confidence * factor;
      normalizedAreas.push({
        area: area,
        confidence: normalizedConf,
        label: getConfidenceLabel(normalizedConf)
      });
    }
    
    normalizedAreas.sort((a, b) => b.confidence - a.confidence);
    
    const grouped = {};
    for (const item of normalizedAreas) {
      if (item.label === 'F') continue;
      if (!grouped[item.label]) grouped[item.label] = [];
      grouped[item.label].push(convertAreaCode(item.area));
    }
    
    return grouped;
  }

  function getConfidenceLabel(confidence) {
    if (confidence >= 0.8) return 'A';
    if (confidence >= 0.6) return 'B';
    if (confidence >= 0.4) return 'C';
    if (confidence >= 0.2) return 'D';
    if (confidence >= 0) return 'E';
    return 'F';
  }

  function convertAreaCode(code) {
    const areaMap = {
      "900": "地域未設定", "901": "地域不明", "905": "日本以外",
      "10": "北海道 石狩", "15": "北海道 渡島", "20": "北海道 檜山", "25": "北海道 後志",
      "30": "北海道 空知", "35": "北海道 上川", "40": "北海道 留萌", "45": "北海道 宗谷",
      "50": "北海道 網走", "55": "北海道 胆振", "60": "北海道 日高", "65": "北海道 十勝",
      "70": "北海道 釧路", "75": "北海道 根室", "100": "青森津軽", "105": "青森三八上北",
      "106": "青森下北", "110": "岩手沿岸北部", "111": "岩手沿岸南部", "115": "岩手内陸",
      "120": "宮城北部", "125": "宮城南部", "130": "秋田沿岸", "135": "秋田内陸",
      "140": "山形庄内", "141": "山形最上", "142": "山形村山", "143": "山形置賜",
      "150": "福島中通り", "151": "福島浜通り", "152": "福島会津", "200": "茨城北部",
      "205": "茨城南部", "210": "栃木北部", "215": "栃木南部", "220": "群馬北部",
      "225": "群馬南部", "230": "埼玉北部", "231": "埼玉南部", "232": "埼玉秩父",
      "240": "千葉北東部", "241": "千葉北西部", "242": "千葉南部", "250": "東京",
      "255": "伊豆諸島北部", "260": "伊豆諸島南部", "265": "小笠原", "270": "神奈川東部",
      "275": "神奈川西部", "300": "新潟上越", "301": "新潟中越", "302": "新潟下越",
      "305": "新潟佐渡", "310": "富山東部", "315": "富山西部", "320": "石川能登",
      "325": "石川加賀", "330": "福井嶺北", "335": "福井嶺南", "340": "山梨東部",
      "345": "山梨中・西部", "350": "長野北部", "351": "長野中部", "355": "長野南部",
      "400": "岐阜飛騨", "405": "岐阜美濃", "410": "静岡伊豆", "411": "静岡東部",
      "415": "静岡中部", "416": "静岡西部", "420": "愛知東部", "425": "愛知西部",
      "430": "三重北中部", "435": "三重南部", "440": "滋賀北部", "445": "滋賀南部",
      "450": "京都北部", "455": "京都南部", "460": "大阪北部", "465": "大阪南部",
      "470": "兵庫北部", "475": "兵庫南部", "480": "奈良", "490": "和歌山北部",
      "495": "和歌山南部", "500": "鳥取東部", "505": "鳥取中・西部", "510": "島根東部",
      "515": "島根西部", "514": "島根隠岐", "520": "岡山北部", "525": "岡山南部",
      "530": "広島北部", "535": "広島南部", "540": "山口北部", "545": "山口中・東部",
      "541": "山口西部", "550": "徳島北部", "555": "徳島南部", "560": "香川",
      "570": "愛媛東予", "575": "愛媛中予", "576": "愛媛南予", "580": "高知東部",
      "581": "高知中部", "582": "高知西部", "600": "福岡福岡", "601": "福岡北九州",
      "602": "福岡筑豊", "605": "福岡筑後", "610": "佐賀北部", "615": "佐賀南部",
      "620": "長崎北部", "625": "長崎南部", "630": "長崎壱岐・対馬", "635": "長崎五島",
      "640": "熊本阿蘇", "641": "熊本熊本", "645": "熊本球磨", "646": "熊本天草・芦北",
      "650": "大分北部", "651": "大分中部", "655": "大分西部", "656": "大分南部",
      "660": "宮崎北部平野部", "661": "宮崎北部山沿い", "665": "宮崎南部平野部", "666": "宮崎南部山沿い",
      "670": "鹿児島薩摩", "675": "鹿児島大隅", "680": "種子島・屋久島", "685": "鹿児島奄美",
      "700": "沖縄本島北部", "701": "沖縄本島中南部", "702": "沖縄久米島", "705": "沖縄八重山",
      "706": "沖縄宮古島", "710": "沖縄大東島"
    };
    return areaMap[code] || code;
  }

  function updateConfidenceDisplay(grouped) {
    const labels = ['A', 'B', 'C', 'D', 'E'];
    let html = '<div class="font-bold col-span-2">各地域の相対的な信頼度</div>';
    
    for (const label of labels) {
      if (grouped[label] && grouped[label].length > 0) {
        html += `<div><span class="text-sm x-confidence x-confidence-${label}">${label}</span></div>`;
        html += `<div class="text-sm py-0.5">${grouped[label].join('、')}</div>`;
      }
    }
    
    html += '<div class="text-xs col-span-2">信頼度は揺れの強さを示すものではありません。「相対的な差」「分布の拡がり」に着目してご覧ください。</div>';
    
    confidenceDisplay.innerHTML = html;
  }

  function preloadImages() {
    if (hasPreloaded || isPreloading) return Promise.resolve();
    
    isPreloading = true;
    playBtn.textContent = '読み込み中...';
    playBtn.disabled = true;
    
    const imagePromises = timeseriesData.map(data => {
      return new Promise((resolve) => {
        const objectId = extractObjectId(data);
        
        if (objectId) {
          const imageUrl = `https://cdn.p2pquake.net/app/web/userquake?id=${objectId}&suffix=_trim`;
          const img = new Image();
          
          img.onload = () => {
            preloadedImages.set(objectId, img);
            console.log('Preloaded image for objectId:', objectId);
            resolve();
          };
          
          img.onerror = () => {
            console.warn('Failed to preload image for objectId:', objectId);
            resolve();
          };
          
          img.src = imageUrl;
        } else {
          resolve();
        }
      });
    });
    
    return Promise.all(imagePromises).then(() => {
      hasPreloaded = true;
      isPreloading = false;
      playBtn.textContent = '▶ 再生';
      playBtn.disabled = false;
    });
  }

  function extractObjectId(data) {
    if (data._id) {
      if (typeof data._id === 'string') {
        return data._id;
      } else if (data._id.$oid) {
        return data._id.$oid;
      } else if (data._id.toString) {
        return data._id.toString();
      }
    }
    return null;
  }

  function updateImage(data) {
    const objectId = extractObjectId(data);
    
    if (objectId) {
      const preloadedImg = preloadedImages.get(objectId);
      if (preloadedImg) {
        timelineImage.src = preloadedImg.src;
        timelineImageLink.href = preloadedImg.src;
      } else {
        console.warn('Preloaded image not found for objectId:', objectId, 'Available keys:', Array.from(preloadedImages.keys()));
        const imageUrl = `https://cdn.p2pquake.net/app/web/userquake?id=${objectId}&suffix=_trim`;
        timelineImage.src = imageUrl;
        timelineImageLink.href = imageUrl;
      }
    }
  }

  playBtn.addEventListener('click', () => {
    if (isPreloading) return;
    
    if (isPlaying) {
      stopAnimation();
    } else {
      if (!hasPreloaded) {
        preloadImages().then(() => {
          if (currentTimeSeconds >= totalDurationSeconds) {
            currentTimeSeconds = 0;
            slider.value = 0;
            updateDisplay();
          }
          startAnimation();
        });
      } else {
        if (currentTimeSeconds >= totalDurationSeconds) {
          currentTimeSeconds = 0;
          slider.value = 0;
          updateDisplay();
        }
        startAnimation();
      }
    }
  });

  slider.addEventListener('input', (e) => {
    stopAnimation();
    currentTimeSeconds = parseInt(e.target.value);
    updateDisplay();
  });

  speedSelect.addEventListener('change', () => {
    if (isPlaying) {
      stopAnimation();
      startAnimation();
    }
  });

  function startAnimation() {
    if (timeseriesData.length <= 1) return;
    
    isPlaying = true;
    playBtn.textContent = '■ 停止';
    playBtn.dataset.playing = 'true';
    
    const speedMultiplier = parseFloat(speedSelect.value);
    const baseInterval = 1000;
    const actualInterval = baseInterval / speedMultiplier;
    
    animationInterval = setInterval(() => {
      if (currentTimeSeconds >= totalDurationSeconds) {
        stopAnimation();
        return;
      }
      
      currentTimeSeconds++;
      slider.value = currentTimeSeconds;
      updateDisplay();
    }, actualInterval);
  }

  function stopAnimation() {
    isPlaying = false;
    playBtn.textContent = '▶ 再生';
    playBtn.dataset.playing = 'false';
    
    if (animationInterval) {
      clearInterval(animationInterval);
      animationInterval = null;
    }
  }
}