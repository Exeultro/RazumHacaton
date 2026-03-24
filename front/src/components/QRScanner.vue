<template>
  <div class="qr-scanner">
    <div id="qr-reader" style="width: 300px"></div>
    <div v-if="scanResult">Результат: {{ scanResult }}</div>
    <el-button @click="startScan" type="primary">Запустить сканер</el-button>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue';
import { Html5Qrcode } from 'html5-qrcode';

const emit = defineEmits(['scan']);
let html5QrCode = null;
const scanResult = ref('');

const startScan = () => {
  if (!html5QrCode) {
    html5QrCode = new Html5Qrcode('qr-reader');
  }
  html5QrCode.start(
    { facingMode: 'environment' },
    { fps: 10, qrbox: { width: 250, height: 250 } },
    (decodedText) => {
      scanResult.value = decodedText;
      emit('scan', decodedText);
      html5QrCode.stop();
    },
    (errorMessage) => {
      console.warn(errorMessage);
    }
  );
};

onUnmounted(() => {
  if (html5QrCode && html5QrCode.isScanning) {
    html5QrCode.stop();
  }
});
</script>