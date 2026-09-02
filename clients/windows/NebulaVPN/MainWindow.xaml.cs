using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Microsoft.UI.Xaml.Media.Animation;
using NebulaVPN.ViewModels;
using System.ComponentModel;

namespace NebulaVPN.Views;

public sealed partial class MainWindow : Window
{
    private readonly MainViewModel _vm;

    public MainWindow()
    {
        InitializeComponent();
        _vm = new MainViewModel();
        _vm.PropertyChanged += OnPropertyChanged;
        ExtendsContentIntoTitleBar = true;

        if (Content is FrameworkElement el)
            el.RequestedTheme = ElementTheme.Dark;

        UpdateServerUI();
    }

    private void OnPropertyChanged(object? sender, PropertyChangedEventArgs e)
    {
        DispatcherQueue.TryEnqueue(() =>
        {
            switch (e.PropertyName)
            {
                case nameof(MainViewModel.IsConnected): UpdateConnectionUI(); break;
                case nameof(MainViewModel.SelectedServer): UpdateServerUI(); break;
                case nameof(MainViewModel.DownloadBytes):
                case nameof(MainViewModel.UploadBytes): UpdateDataUI(); break;
            }
        });
    }

    private void UpdateConnectionUI()
    {
        if (_vm.IsConnected)
        {
            StatusText.Text = "Connected";
            ConnectButton.Content = "Disconnect";
            ConnectButton.Style = (Style)Resources["DangerButton"];
            DataBar.Visibility = Visibility.Visible;

            var anim = new DoubleAnimation
            {
                From = 0, To = 1,
                Duration = new Duration(TimeSpan.FromSeconds(0.6)),
                EasingFunction = new CubicEase { EasingMode = EasingMode.EaseOut }
            };
            Storyboard.SetTarget(anim, RingActive);
            Storyboard.SetTargetProperty(anim, "Opacity");
            new Storyboard { Children = { anim } }.Begin();
        }
        else
        {
            StatusText.Text = "Not Connected";
            DurationText.Text = "";
            ConnectButton.Content = "Connect";
            ConnectButton.Style = (Style)Resources["AccentButtonStyle"];
            DataBar.Visibility = Visibility.Collapsed;
            RingActive.Opacity = 0;
        }
    }

    private void UpdateServerUI()
    {
        if (_vm.SelectedServer != null)
        {
            ServerFlag.Text = _vm.SelectedServer.CountryFlag;
            ServerCity.Text = _vm.SelectedServer.City;
            ServerDetail.Text = $"{_vm.SelectedServer.LatencyMs}ms · {_vm.SelectedServer.LoadPercent}% load";
        }
    }

    private void UpdateDataUI()
    {
        DownloadText.Text = FormatBytes(_vm.DownloadBytes);
        UploadText.Text = FormatBytes(_vm.UploadBytes);
    }

    private static string FormatBytes(long bytes)
    {
        string[] suff = { "B", "KB", "MB", "GB", "TB" };
        int i = 0;
        double d = bytes;
        while (d >= 1024 && i < suff.Length - 1) { d /= 1024; i++; }
        return $"{d:F1} {suff[i]}";
    }

    private void OnConnectClick(object _, RoutedEventArgs _) => _vm.ToggleConnection();
    private void OnServerClick(object _, RoutedEventArgs _) { /* Show picker dialog */ }
}
