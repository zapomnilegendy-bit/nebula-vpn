using System.ComponentModel;
using System.Runtime.CompilerServices;
using System.Timers;
using Windows.Storage;

namespace NebulaVPN.ViewModels;

public class MainViewModel : INotifyPropertyChanged
{
    private bool _isConnected;
    private ServerItem? _selectedServer;
    private long _downloadBytes;
    private long _uploadBytes;
    private int _connectionSeconds;
    private System.Timers.Timer? _trafficTimer;
    private System.Timers.Timer? _durationTimer;

    // Settings (stored in ApplicationData)
    public bool KillSwitchEnabled { get; set; } = true;
    public bool AdaptiveProtocol { get; set; } = true;
    public bool LaunchAtStartup { get; set; } = false;

    public bool IsConnected
    {
        get => _isConnected;
        private set { _isConnected = value; OnPropertyChanged(); }
    }

    public ServerItem? SelectedServer
    {
        get => _selectedServer;
        set { _selectedServer = value; OnPropertyChanged(); UpdateServerUI(); }
    }

    public long DownloadBytes
    {
        get => _downloadBytes;
        private set { _downloadBytes = value; OnPropertyChanged(); }
    }

    public long UploadBytes
    {
        get => _uploadBytes;
        private set { _uploadBytes = value; OnPropertyChanged(); }
    }

    public int ConnectionSeconds
    {
        get => _connectionSeconds;
        set { _connectionSeconds = value; OnPropertyChanged(); OnPropertyChanged(nameof(FormattedDuration)); }
    }

    public string FormattedDuration
    {
        get
        {
            var ts = TimeSpan.FromSeconds(ConnectionSeconds);
            return ts.TotalHours >= 1
                ? $"{(int)ts.TotalHours}h {ts.Minutes}m"
                : $"{ts.Minutes}m {ts.Seconds}s";
        }
    }

    public List<ServerItem> Servers { get; } = new()
    {
        new("auto", "Fastest Server", "⚡", "", 0, 0, false),
        new("eu-fra-01", "Frankfurt", "🇩🇪", "fra1.nebula", 12, 42, false),
        new("eu-ams-01", "Amsterdam", "🇳🇱", "ams1.nebula", 18, 58, false),
        new("gb-lon-01", "London", "🇬🇧", "lon1.nebula", 24, 55, false),
        new("us-nyc-01", "New York", "🇺🇸", "nyc1.nebula", 89, 71, true),
        new("us-sfo-01", "San Francisco", "🇺🇸", "sfo1.nebula", 105, 34, false),
        new("jp-tokyo-01", "Tokyo", "🇯🇵", "tky1.nebula", 185, 23, true),
        new("au-syd-01", "Sydney", "🇦🇺", "syd1.nebula", 280, 12, true),
    };

    public void ToggleConnection()
    {
        if (IsConnected) Disconnect(); else Connect();
    }

    public void Connect()
    {
        if (SelectedServer == null) return;
        IsConnected = true;
        StartTimers();
    }

    public void Disconnect()
    {
        IsConnected = false;
        StopTimers();
        DownloadBytes = 0;
        UploadBytes = 0;
        ConnectionSeconds = 0;
    }

    private void StartTimers()
    {
        _trafficTimer = new(2000);
        _trafficTimer.Elapsed += (_, _) =>
        {
            DownloadBytes += Random.Shared.Next(1000, 50000);
            UploadBytes += Random.Shared.Next(100, 5000);
        };
        _trafficTimer.AutoReset = true;
        _trafficTimer.Start();

        _durationTimer = new(1000);
        _durationTimer.Elapsed += (_, _) => ConnectionSeconds++;
        _durationTimer.AutoReset = true;
        _durationTimer.Start();
    }

    private void StopTimers()
    {
        _trafficTimer?.Stop(); _trafficTimer?.Dispose(); _trafficTimer = null;
        _durationTimer?.Stop(); _durationTimer?.Dispose(); _durationTimer = null;
    }

    private void UpdateServerUI() => OnPropertyChanged(nameof(SelectedServer));

    // DPAPI secure storage
    public void SaveCredentials(string token)
    {
        var vault = new Windows.Security.Credentials.PasswordVault();
        vault.Add(new("NebulaVPN", "auth_token", token));
    }

    public string? LoadCredentials()
    {
        try
        {
            var vault = new Windows.Security.Credentials.PasswordVault();
            var cred = vault.Retrieve("NebulaVPN", "auth_token");
            cred.RetrievePassword();
            return cred.Password;
        }
        catch { return null; }
    }

    public event PropertyChangedEventHandler? PropertyChanged;
    private void OnPropertyChanged([CallerMemberName] string? n = null)
        => PropertyChanged?.Invoke(this, new(n));
}

public record ServerItem(
    string Id,
    string City,
    string CountryFlag,
    string DisplayName,
    int LatencyMs,
    int LoadPercent,
    bool IsPremium
);
