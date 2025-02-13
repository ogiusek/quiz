import { createContext, useEffect, useState } from "react";

interface Theme {
  Toggle(): void
  Dark(): void
  Light(): void
  State(): 'dark' | 'light'
}

export const ThemeContext = createContext<Theme>({
  Toggle() { throw new Error("not implemented") },
  Dark() { throw new Error("not implemented") },
  Light() { throw new Error("not implemented") },
  State() { throw new Error("not implemented") },
})

export const ThemeService: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');
  const classList = document.querySelector('body')?.classList

  useEffect(() => {
    if (theme == 'dark') classList?.add('dark')
    if (theme == 'light') classList?.remove('dark')
  })

  return <ThemeContext.Provider value={{
    Toggle() { setTheme(theme == 'dark' ? 'light' : 'dark') },
    Dark() { setTheme('dark') },
    Light() { setTheme('light') },
    State() { return theme },
  }}>
    {children}
  </ThemeContext.Provider>
}