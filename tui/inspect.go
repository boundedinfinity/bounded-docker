package tui

import (
	"fmt"
	"sort"
	"strconv"

	moby "github.com/moby/moby/client"
	"github.com/rivo/tview"
)

func newInspect(tui *tui, id string) *inspect {
	inspect := &inspect{
		tui: tui, containerId: id,
	}

	inspect.view = tview.NewTable().
		SetBorders(true).
		SetFixed(1, 3).
		SetEvaluateAllRows(true).
		SetSelectable(true, false)
	inspect.view.SetBorder(true).SetTitle(" [ Inspect: " + id + " ] ")

	return inspect
}

type inspect struct {
	tui         *tui
	view        *tview.Table
	rows        [][]string
	containerId string
}

func (this *inspect) Data() *tview.Table {
	return this.view
}

// Value returns the value of the currently selected row, skipping the header.
func (this *inspect) Value() (string, bool) {
	row, _ := this.view.GetSelection()

	if row <= 0 || row >= len(this.rows) {
		return "", false
	}

	value := this.rows[row][1]
	return value, value != ""
}

func (this *inspect) Set(result moby.ContainerInspectResult) {
	container := result.Container
	rows := [][]string{}
	index := 0

	add := func(key string, values ...string) {
		rows, index = Utils.Tview.multiRowi(rows, index, key, values...)
	}

	add("Id", container.ID)
	add("Name", container.Name)
	add("Created", container.Created)
	add("Image", container.Image)
	add("Platform", container.Platform)
	add("Driver", container.Driver)
	add("Path", container.Path)
	add("Args", container.Args...)
	add("RestartCount", strconv.Itoa(container.RestartCount))
	add("MountLabel", container.MountLabel)
	add("ProcessLabel", container.ProcessLabel)
	add("AppArmorProfile", container.AppArmorProfile)
	add("ResolvConfPath", container.ResolvConfPath)
	add("HostnamePath", container.HostnamePath)
	add("HostsPath", container.HostsPath)
	add("LogPath", container.LogPath)
	add("ExecIDs", container.ExecIDs...)

	// Only populated when ContainerInspectOptions.Size is set.
	if container.SizeRw != nil {
		add("SizeRw", Utils.size2Str(*container.SizeRw))
	}

	if container.SizeRootFs != nil {
		add("SizeRootFs", Utils.size2Str(*container.SizeRootFs))
	}

	if state := container.State; state != nil {
		add("State.Status", string(state.Status))
		add("State.Running", strconv.FormatBool(state.Running))
		add("State.Paused", strconv.FormatBool(state.Paused))
		add("State.Restarting", strconv.FormatBool(state.Restarting))
		add("State.OOMKilled", strconv.FormatBool(state.OOMKilled))
		add("State.Dead", strconv.FormatBool(state.Dead))
		add("State.Pid", strconv.Itoa(state.Pid))
		add("State.ExitCode", strconv.Itoa(state.ExitCode))
		add("State.Error", state.Error)
		add("State.StartedAt", state.StartedAt)
		add("State.FinishedAt", state.FinishedAt)

		if health := state.Health; health != nil {
			add("State.Health.Status", string(health.Status))
			add("State.Health.FailingStreak", strconv.Itoa(health.FailingStreak))
		}
	}

	if config := container.Config; config != nil {
		add("Config.Hostname", config.Hostname)
		add("Config.Domainname", config.Domainname)
		add("Config.User", config.User)
		add("Config.Image", config.Image)
		add("Config.WorkingDir", config.WorkingDir)
		add("Config.Tty", strconv.FormatBool(config.Tty))
		add("Config.OpenStdin", strconv.FormatBool(config.OpenStdin))
		add("Config.StdinOnce", strconv.FormatBool(config.StdinOnce))
		add("Config.AttachStdin", strconv.FormatBool(config.AttachStdin))
		add("Config.AttachStdout", strconv.FormatBool(config.AttachStdout))
		add("Config.AttachStderr", strconv.FormatBool(config.AttachStderr))
		add("Config.NetworkDisabled", strconv.FormatBool(config.NetworkDisabled))
		add("Config.StopSignal", config.StopSignal)
		add("Config.Env", config.Env...)
		add("Config.Cmd", config.Cmd...)
		add("Config.Entrypoint", config.Entrypoint...)
		add("Config.Shell", config.Shell...)
		add("Config.OnBuild", config.OnBuild...)
		add("Config.Labels", Utils.Docker.map2Lines(config.Labels)...)
		add("Config.Volumes", Utils.Docker.volumes2Lines(config.Volumes)...)
		add("Config.ExposedPorts", Utils.Docker.exposedPorts2Lines(config.ExposedPorts)...)

		if config.StopTimeout != nil {
			add("Config.StopTimeout", strconv.Itoa(*config.StopTimeout))
		}
	}

	if host := container.HostConfig; host != nil {
		add("HostConfig.NetworkMode", string(host.NetworkMode))
		add("HostConfig.RestartPolicy", string(host.RestartPolicy.Name))
		add("HostConfig.Privileged", strconv.FormatBool(host.Privileged))
		add("HostConfig.AutoRemove", strconv.FormatBool(host.AutoRemove))
		add("HostConfig.ReadonlyRootfs", strconv.FormatBool(host.ReadonlyRootfs))
		add("HostConfig.Binds", host.Binds...)
		add("HostConfig.CapAdd", host.CapAdd...)
		add("HostConfig.CapDrop", host.CapDrop...)
		add("HostConfig.Dns", Utils.Docker.addrs2Lines(host.DNS)...)
		add("HostConfig.Memory", Utils.size2Str(host.Memory))
		add("HostConfig.CpuShares", strconv.FormatInt(host.CPUShares, 10))
	}

	for i, mount := range container.Mounts {
		prefix := fmt.Sprintf("Mounts[%d]", i)
		add(prefix+".Type", string(mount.Type))
		add(prefix+".Name", mount.Name)
		add(prefix+".Source", mount.Source)
		add(prefix+".Destination", mount.Destination)
		add(prefix+".Driver", mount.Driver)
		add(prefix+".Mode", mount.Mode)
		add(prefix+".RW", strconv.FormatBool(mount.RW))
		add(prefix+".Propagation", string(mount.Propagation))
	}

	if settings := container.NetworkSettings; settings != nil {
		add("NetworkSettings.SandboxID", settings.SandboxID)
		add("NetworkSettings.SandboxKey", settings.SandboxKey)
		add("NetworkSettings.Ports", Utils.Docker.portMap2Lines(settings.Ports)...)

		names := make([]string, 0, len(settings.Networks))

		for name := range settings.Networks {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			endpoint := settings.Networks[name]
			if endpoint == nil {
				continue
			}

			prefix := "NetworkSettings.Networks." + name
			add(prefix+".NetworkID", endpoint.NetworkID)
			add(prefix+".EndpointID", endpoint.EndpointID)
			add(prefix+".Gateway", endpoint.Gateway.String())
			add(prefix+".IPAddress", endpoint.IPAddress.String())
			add(prefix+".MacAddress", endpoint.MacAddress.String())
			add(prefix+".Aliases", endpoint.Aliases...)
		}
	}

	if driver := container.GraphDriver; driver != nil {
		add("GraphDriver.Name", driver.Name)
		add("GraphDriver.Data", Utils.Docker.map2Lines(driver.Data)...)
	}

	this.tui.queueDraw(func() {
		this.view.Clear()
		this.rows = rows

		for r, row := range rows {
			for key, val := range row {
				cell := tview.NewTableCell(tview.Escape(val))

				if r == 0 {
					cell.SetSelectable(false)
				}

				cell.SetExpansion(1)
				this.view.SetCell(r, key, cell)
			}
		}

		this.view.ScrollToBeginning()
	})
}
