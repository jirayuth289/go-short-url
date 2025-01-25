import React, { useState, useMemo } from 'react'
import axios from 'axios'
import './App.css'
import { debounce } from './utils/helper'

function App() {
  const [url, setURL] = useState('')
  const [shortURL, setShortURL] = useState('')

  const handleSubmitURL = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    try {
      const response = await axios.post('http://localhost:8080/shorten', { url })
      const { shortURL } = response.data;
      setShortURL(shortURL)
    } catch (error) {
      console.error(error)
    }
  }

  const handleOpenURL = () => {
    if (shortURL) {
      window.open(`http://localhost:8080/${shortURL}`, '_blank')
    }
  }

  const handleCopyURL = () => {
    if (shortURL) {
      navigator.clipboard.writeText(`http://localhost:8080/${shortURL}`)
    }
  }

  const handleInputURLChange = (e: React.BaseSyntheticEvent) => {
    setURL(e.target.value)
  }

  const debouncedInputHandler = useMemo(() => debounce(handleInputURLChange, 1000), [])

  return (
    <div>
      <form onSubmit={handleSubmitURL}>
        <h1>Shorten URL</h1>
        <input type='text' placeholder='Your Link' onChange={debouncedInputHandler} />
        <button type="submit">Submit</button>
      </form>
      {shortURL && (
        <div>
          <p>Shorten URL: <a href={`http://localhost:8080/${shortURL}`} target='_blank' rel='noopener noreferrer' />{`http://localhost:8080/${shortURL}`}</p>
          <button type='button' onClick={handleOpenURL}>Open URL</button>
          <button type='button' onClick={handleCopyURL}>Copy URL</button>
        </div>
      )}
    </div>
  )
}

export default App
