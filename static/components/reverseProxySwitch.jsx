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