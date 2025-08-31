import { useEffect, useState } from "react"
import { loadConfig, type AppConfig } from "../services/env"
import { initLiff, getGroupIdOrThrow, getProfileSafe } from "../services/liff"
import { postScan } from "../api/client"
import { mark, getTimings } from "../ui/timing"
import { showToast } from "../ui/dom"

function App() {
  const [cfg, setCfg] = useState<AppConfig | null>(null)
  const [groupId, setGroupId] = useState("")
  const [profile, setProfile] = useState({ displayName: "", userId: "" })
  const [error, setError] = useState("")
  const [ready, setReady] = useState(false)
  const [isScanning, setIsScanning] = useState(false)
  const [lastScanResult, setLastScanResult] = useState<string | null>(null)

  useEffect(() => {
    ;(async () => {
      try {
        const c = await loadConfig()
        setCfg(c)
        await initLiff(c)
        const gid = await getGroupIdOrThrow()
        setGroupId(gid)
        const p = await getProfileSafe()
        setProfile(p)
        setReady(true)
      } catch (e: any) {
        setError(e.message ?? String(e))
      }
    })()
  }, [])

  const handleScan = async () => {
    if (!cfg) return

    setIsScanning(true)
    try {
      mark("t0")
      await postScan(cfg, { groupId, qrText: "dummy", ...profile })
      mark("t2")

      // Simulate scan result for demo
      setLastScanResult("https://example.com/scanned-content")
      showToast("scanned")

      if (cfg!.env !== "prod") {
        const debug = document.getElementById("debug")
        if (debug) debug.textContent = JSON.stringify(getTimings())
      }
    } catch (e) {
      console.error("Scan failed:", e)
    } finally {
      setIsScanning(false)
    }
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6 bg-background">
        <div className="bg-card rounded-3xl p-8 max-w-sm w-full text-center shadow-lg">
          <div className="w-16 h-16 bg-destructive/10 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-destructive" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
              />
            </svg>
          </div>
          <h2 className="text-xl font-semibold text-card-foreground mb-2">Oops! Something went wrong</h2>
          <p className="text-muted-foreground text-sm">{error}</p>
        </div>
      </div>
    )
  }

  if (!ready) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6 bg-background">
        <div className="bg-card rounded-3xl p-8 max-w-sm w-full text-center shadow-lg">
          <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4 animate-pulse-gentle">
            <svg className="w-8 h-8 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4"
              />
            </svg>
          </div>
          <h2 className="text-xl font-semibold text-card-foreground mb-2">Getting ready...</h2>
          <p className="text-muted-foreground text-sm">Setting up your QR scanner</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="bg-emerald-50/70 border-b border-emerald-100">
        <div className="max-w-sm mx-auto px-6 py-4">
          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-2xl font-semibold text-emerald-600">QR Scanner</h1>
              <p className="text-sm text-emerald-700/80">Hello, {profile.displayName || "Demo User"}! 👋</p>
            </div>
            <div className="w-10 h-10 bg-emerald-100 rounded-full flex items-center justify-center">
              {/* Camera icon */}
              <svg className="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8h4l2-3h6l2 3h4v11a2 2 0 01-2 2H5a2 2 0 01-2-2V8z" />
                <circle cx="12" cy="14" r="3" strokeWidth={2} />
              </svg>
            </div>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-sm mx-auto px-6 py-8">
        {/* Big circular scan button */}
        <button
          onClick={handleScan}
          className={`mx-auto block w-40 h-40 rounded-full bg-emerald-500 text-white shadow-lg shadow-emerald-300/30 flex items-center justify-center ${
            isScanning ? "animate-pulse" : "hover:scale-105 transition-transform"
          }`}
        >
          {/* QR icon */}
          <svg className="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5h4v4H3zM17 5h4v4h-4zM3 15h4v4H3zM13 13h2v2h-2zM11 17h2v2h-2zM15 17h2v2h-2zM17 13h4v2h-2v2h-2v-4zM7 7h6M7 17h2M13 7h2" />
          </svg>
        </button>

        {/* Title and subtitle */}
        <div className="mt-8 text-center">
          <h2 className="text-3xl font-bold text-emerald-600">{isScanning ? "Scanning..." : "Tap to Scan"}</h2>
          <p className="text-slate-500 mt-2">Point your camera at a QR code</p>
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-2 gap-4 mt-8">
          <button className="bg-emerald-50 border border-emerald-100 rounded-2xl p-6 text-center shadow-sm hover:bg-emerald-50/80">
            <div className="w-9 h-9 bg-emerald-100 rounded-full flex items-center justify-center mx-auto mb-3">
              {/* clock icon */}
              <svg className="w-4.5 h-4.5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <span className="text-emerald-600 font-medium">History</span>
          </button>

          <button className="bg-emerald-50 border border-emerald-100 rounded-2xl p-6 text-center shadow-sm hover:bg-emerald-50/80">
            <div className="w-9 h-9 bg-emerald-100 rounded-full flex items-center justify-center mx-auto mb-3">
              {/* settings gear */}
              <svg className="w-5 h-5 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <span className="text-emerald-600 font-medium">Settings</span>
          </button>
        </div>

        {/* Debug area */}
        {cfg?.env !== "prod" && (
          <div className="bg-emerald-50 border border-emerald-100 rounded-2xl p-4 mt-8 text-slate-600">
            <div className="text-sm font-medium mb-2">Debug Info:</div>
            <pre id="debug" className="text-xs text-slate-500"></pre>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
