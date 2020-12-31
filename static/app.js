class SwitchContainer extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: null,
            items: []
        };
        this.reloadReverseProxies = this.reloadReverseProxies.bind(this);
    }

    reloadReverseProxies() {
        fetch("/serverList")
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        items: result
                    });
                }, (error) => {
                    this.setState({
                        error
                    });
                }
            )
    }
    componentDidMount() {
        this.reloadReverseProxies()
    }
    render() {
        const error = this.state.error;
        const titleStyle = {
            margin:"50px 0 0 200px"
        };
        const reloadButtonStyle = {
            margin:"50px 0 50px 300px"
        };
        const proxyContainerStyle = {
            margin:"0 0 0 120px"
        };
        if (!error) {
            return (
                <div>
                    <h1 style={titleStyle}>Reverse Proxy Switcher</h1>
                    <button style={reloadButtonStyle} className="btn btn-warning" onClick={this.reloadReverseProxies}
                            type="button">Reload
                    </button>
                    <ul style={proxyContainerStyle}>
                        {this.state.items.map((item, index) => (
                            <ReverseProxySwitch name={item.Name} reloadHandler={this.reloadReverseProxies} currentValue={item.CurrentValue} isLocalRunning={item.IsLocalRunning} key={index}/>
                        ))}
                    </ul>
                </div>
            );
        } else {
            return (
                <div>
                    <h1>Error!</h1>
                </div>
            )
        }
    }
}
class ReverseProxySwitch extends React.Component {

    constructor(props){
        super(props);
        this.switchServer = this.switchServer.bind(this);
    }

    switchServer() {
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        };
        fetch('/switchReverseProxyServer?serverName=' + this.props.name, requestOptions)
            .then(data => this.props.reloadHandler());
    }

    render() {
        const componentStyle = {
            margin:"10px 0 0 0"
        };

        const titleStyle = {
            margin:"5px 20px 0 0"
        };

        const localButtonStyle = {
            margin:"0 10px 0 0"
        };
        return (

            <div style={componentStyle} className="input-group">
                <span style={titleStyle}>{this.props.name} ({this.props.currentValue === 1 ? "Remote" : "Local"})</span>
                <div className="btn-group" id="button-addon4">
                    <button disabled={this.props.currentValue !== 1 || !this.props.isLocalRunning} style={localButtonStyle} className="btn btn-primary" onClick={this.switchServer}
                            type="button">Local {!this.props.isLocalRunning? "(Not Working)" : ""}
                    </button>

                    <button  disabled={this.props.currentValue === 1} className="btn btn-success"  onClick={this.switchServer}
                             id="switchUrl" type="button">Remote Server
                    </button>
                </div>
            </div>
        )
    }
}

ReactDOM.render(
    <SwitchContainer />,
    document.getElementById("root")
);