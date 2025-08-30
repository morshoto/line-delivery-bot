import { useEffect, useState } from 'react';
import { loadConfig, type AppConfig } from '../services/env';
import { initLiff, getGroupIdOrThrow, getProfileSafe } from '../services/liff';
import { postScan } from '../api/client';
import { mark, getTimings } from '../ui/timing';
import { showToast } from '../ui/dom';

function App() {
  const [cfg, setCfg] = useState<AppConfig | null>(null);
  const [groupId, setGroupId] = useState('');
  const [profile, setProfile] = useState({ displayName: '', userId: '' });
  const [error, setError] = useState('');
  const [ready, setReady] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const c = await loadConfig();
        setCfg(c);
        await initLiff(c);
        const gid = await getGroupIdOrThrow();
        setGroupId(gid);
        const p = await getProfileSafe();
        setProfile(p);
        setReady(true);
      } catch (e: any) {
        setError(e.message ?? String(e));
      }
    })();
  }, []);

  const handleScan = async () => {
    if (!cfg) return;
    mark('t0');
    await postScan(cfg, { groupId, qrText: 'dummy', ...profile });
    mark('t2');
    showToast('scanned');
    if (cfg!.env !== 'prod') {
      const debug = document.getElementById('debug');
      if (debug) debug.textContent = JSON.stringify(getTimings());
    }
  };

  if (error) return <p>{error}</p>;
  if (!ready) return <p>Loading...</p>;
  return (
      <div>
        <h1>QR Scanner</h1>
        <button onClick={handleScan}>Scan</button>
        {cfg?.env !== 'prod' && <pre id="debug"></pre>}
      </div>
  );
}

export default App;
