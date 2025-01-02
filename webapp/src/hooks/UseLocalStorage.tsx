import { useEffect, useState } from 'react'
import { isLocalStorageAvailable } from '../utils/Storage';

export function useLocalStorage<T>(key: string, defaultValue: T): [T, (value: T) => void] {
  const [value, setValue] = useState<T>(() => {
    return getStorageValue<T>(key, defaultValue);
  });
  useEffect(() => {
    if(isLocalStorageAvailable()) {
      localStorage.setItem(key, JSON.stringify(value));
    }
  }, [key, value]);
  return [value, setValue]
};

function getStorageValue<T>(key: string, defaultValue: T): T {
  if(isLocalStorageAvailable()) {
    const saved: string | null = localStorage.getItem(key);
    if(saved == null) {
      return defaultValue
    }
    const initial = JSON.parse(saved);
    return initial;
  }else {
    return defaultValue;
  }
} 