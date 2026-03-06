<script setup>
import { ref, computed } from 'vue'
import axios from 'axios'

const longUrl = ref('')
const customUrl = ref('')
const shortUrl = ref('')
const loading = ref(false)
const error = ref('')
const copySuccess = ref(false)
const existsLong = ref(false)
const useCustomExpiry = ref(false)
const expiryDate = ref('')

// Determine the base URL for the short links
// In production this would be the actual domain.
// For dev, it's the backend address.
const BASE_URL = 'http://localhost:8082/'

const fullShortUrl = computed(() => {
  if (!shortUrl.value) return ''
  // If the backend returns a full URL, use it. Otherwise append to base.
  if (shortUrl.value.startsWith('http')) return shortUrl.value
  return `${BASE_URL}s/${shortUrl.value}`
})

const validateUrl = (url) => {
  try {
    new URL(url)
    return true
  } catch (_) {
    return false
  }
}

const generateShortUrl = async () => {
  error.value = ''
  shortUrl.value = ''
  copySuccess.value = false
  existsLong.value = false

  if (!longUrl.value) {
    error.value = '请输入长链接'
    return
  }
  
  let urlToSubmit = longUrl.value.trim()
  if (!/^https?:\/\//i.test(urlToSubmit)) {
      urlToSubmit = 'http://' + urlToSubmit
  }

  if (!validateUrl(urlToSubmit)) {
    error.value = '请输入有效的URL格式 (例如: https://example.com)'
    return
  }

  loading.value = true
  try {
    const payload = {
      long_url: urlToSubmit,
      custom_url: customUrl.value || undefined,
      expired_at: useCustomExpiry.value && expiryDate.value ? new Date(expiryDate.value).toISOString() : null
    }
    
    console.log('Sending payload:', JSON.stringify(payload, null, 2)) // Debug log

    const response = await axios.post('http://localhost:8082/api/v1/url/create', payload)

    const errMsg = response.data?.errorMsg
    const isLongUrlExists = /传入的长链接已经存在|ERROR_LONGURL_AVAILABLE/i.test(errMsg)
    const isShortCodeExists = /CustomCode\s*already\s*exists|短码已存在|短链已存在|duplicate|already\s*exists/i.test(errMsg)

    if (isShortCodeExists) {
      error.value = '短码已存在，请更换或留空'
      shortUrl.value = ''
      return
    }

    if (response.data && response.data.rep && response.data.rep.short_url) {
      shortUrl.value = response.data.rep.short_url
      if (isLongUrlExists) {
        existsLong.value = true
        error.value = '传入的长链接已经存在'
      }
    } else {
      error.value = '生成失败：' + (errMsg || '未知错误')
    }
  } catch (err) {
    console.error(err)
    const errMsg = err.response?.data?.errormsg || err.response?.data?.msg
    if (errMsg && /CustomCode\s*already\s*exists|短码已存在|短链已存在|duplicate|already\s*exists/i.test(errMsg)) {
      error.value = '短码已存在，请更换或留空'
      shortUrl.value = ''
    } else if (errMsg && /传入的长链接已经存在|ERROR_LONGURL_AVAILABLE/i.test(errMsg)) {
      existsLong.value = true
      // 如果后端在错误响应中也返回了 rep.short_url，则尽量展示
      const s = err.response?.data?.rep?.short_url
      if (s) shortUrl.value = s
    } else {
      error.value = errMsg || '请求失败，请检查网络或后端服务'
    }
  } finally {
    loading.value = false
  }
}

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(fullShortUrl.value)
    copySuccess.value = true
    setTimeout(() => copySuccess.value = false, 2000)
  } catch (err) {
    console.error('Failed to copy', err)
    error.value = '复制失败，请手动复制'
  }
}
</script>

<template>
  <div class="container">
    <div class="card">
      <div class="header">
        <h1 class="title">🔗 短链接生成器</h1>
        <p class="subtitle">让分享变得简单高效</p>
      </div>
      
      <div class="form-group">
        <div class="input-wrapper">
          <input 
            v-model="longUrl" 
            type="text" 
            placeholder="在此粘贴长链接 (https://...)" 
            :class="{ 'error-border': error }"
            @keyup.enter="generateShortUrl"
          />
        </div>
        
        <div class="input-wrapper">
          <input 
            v-model="customUrl" 
            type="text" 
            placeholder="自定义短码 (可选)" 
            class="custom-input"
          />
        </div>

        <div class="expiry-section">
          <div class="checkbox-wrapper">
            <input 
              type="checkbox" 
              id="expiry-toggle" 
              v-model="useCustomExpiry"
            >
            <label for="expiry-toggle">启用自定义过期时间</label>
          </div>
          
          <transition name="fade">
            <div v-if="useCustomExpiry" class="input-wrapper date-picker-wrapper">
              <label class="date-label">过期时间：</label>
              <input 
                v-model="expiryDate" 
                type="datetime-local" 
                step="1"
                class="date-input"
              />
            </div>
          </transition>
        </div>

        <transition name="fade">
          <p v-if="error" class="error-msg">⚠️ {{ error }}</p>
        </transition>

        <button @click="generateShortUrl" :disabled="loading" class="generate-btn">
          <span v-if="loading" class="loader"></span>
          <span v-else>立即生成</span>
        </button>
      </div>

      <transition name="slide-up">
        <div v-if="shortUrl" class="result-area">
          <p class="success-label">{{ existsLong ? 'ℹ️ 长链接已经存在' : '✅ 生成成功！' }}</p>
          <div class="url-box">
            <a :href="fullShortUrl" target="_blank" class="short-link">{{ fullShortUrl }}</a>
            <button @click="copyToClipboard" class="copy-btn" :class="{ 'copied': copySuccess }">
              {{ copySuccess ? '已复制' : '复制链接' }}
            </button>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<style scoped>
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
}

.card {
  background: rgba(30, 41, 59, 0.7);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 3rem;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  transition: transform 0.3s ease;
}

.card:hover {
  transform: translateY(-5px);
}

.header {
  margin-bottom: 2.5rem;
}

.title {
  font-size: 2.2rem;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(to right, #60a5fa, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  color: #94a3b8;
  margin-top: 0.5rem;
  font-size: 1.1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.input-wrapper input {
  width: 100%;
  padding: 1rem;
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 8px;
  color: white;
  font-size: 1rem;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.input-wrapper input:focus {
  outline: none;
  border-color: #60a5fa;
  box-shadow: 0 0 0 3px rgba(96, 165, 250, 0.2);
}

.error-border {
  border-color: #ef4444 !important;
}

.error-msg {
  color: #ef4444;
  font-size: 0.9rem;
  margin: 0;
  text-align: left;
}

.generate-btn {
  margin-top: 1rem;
  padding: 1rem;
  background: linear-gradient(to right, #3b82f6, #8b5cf6);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s, transform 0.1s;
}

.generate-btn:hover {
  opacity: 0.9;
}

.generate-btn:active {
  transform: scale(0.98);
}

.generate-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.result-area {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.success-label {
  color: #4ade80;
  font-weight: 600;
  margin-bottom: 1rem;
}

.url-box {
  display: flex;
  align-items: center;
  background: rgba(15, 23, 42, 0.8);
  padding: 0.75rem;
  border-radius: 8px;
  border: 1px solid rgba(148, 163, 184, 0.2);
}

.short-link {
  flex: 1;
  color: #60a5fa;
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 1rem;
  font-family: monospace;
  font-size: 1.1rem;
}

.short-link:hover {
  text-decoration: underline;
}

.copy-btn {
  background: #334155;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.copy-btn:hover {
  background: #475569;
}

.copy-btn.copied {
  background: #22c55e;
}

/* Animations */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.4s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

@media (max-width: 640px) {
  .card {
    padding: 1.5rem;
  }
  .title {
    font-size: 1.8rem;
  }
}
</style>
