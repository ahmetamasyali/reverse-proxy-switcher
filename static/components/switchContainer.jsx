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