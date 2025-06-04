import axios from 'axios';

export async function proxyHealthcheck() {
    await axios.get('/api/v1/healthcheck')
        .then((res) => console.log('✅ Proxy success:', res.data))
        .catch((err) => console.error('❌ Proxy error:', err));
}
