import React, { Component } from 'react'

class About extends Component {

  componentDidMount(){
    document.title = 'About - WPdirectory'
  }

  render() {
      return (
        <div className="page page-about">
          <div className="grid-x gutter-x gutter-y small-1">
            <div className="panel">
              <h2>About</h2>
              <p>WPdirectory mirrors the WordPress Plugin and Theme Directories, allowing lightning fast regex search. It does this by using trigram indexing, allowing it to only search relevant files. The projects aims to significantly reduce the time required for searching the Plugin and Theme Directories, saving time and encouraging more use.</p>
              <p>WPdir was inspired by the various Directory Slurper projects by <a href="https://github.com/markjaquith/WordPress-Plugin-Directory-Slurper" rel="noopener noreferrer" target="_blank">Mark Jaquith</a>, <a href="https://github.com/Ipstenu/WordPress-Theme-Directory-Slurper" rel="noopener noreferrer" target="_blank">Ipstenu</a> and many others. It is built using Go on the backend and React on the frontend, it uses Google's <a href="https://github.com/google/codesearch/" rel="noopener noreferrer" target="_blank">codesearch</a> tool for indexing and search.</p>
              <p>If you have feedback, questions or issues please let me know on <a href="https://github.com/wpdirectory/wpdir">github.com/wpdirectory/wpdir</a> or send me a message on the WordPress Slack (peterbooker).</p>
            </div>
            <div className="panel">
              <h2>Privacy</h2>
              <p>WPdirectory does not store any information relating to visistors and/or users. The only information collected is internal metrics, via Prometheus, things like how many searches occured. That is right- no cookies, ads, visitor tracking, etc.</p>
            </div>
          </div>
        </div>
      )
  }
}

export default About