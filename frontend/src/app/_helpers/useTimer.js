import { useState, useEffect, useRef } from "react";

export default function (callback, delay, period) {
  const latestCallback = useRef(() => { });
  const [start, setStart] = useState(false);

  useEffect(() => {
    latestCallback.current = callback;
  });

  useEffect(() => {
    setStart(false);

    setTimeout(() => {
      latestCallback.current();
      setStart(true);
    }, delay);

  }, [delay, period]);

  useEffect(() => {
    if (start && delay !== null) {
      let id = setInterval(() => latestCallback.current(), period)
      return () => clearInterval(id);
    }
  }, [start, delay, period]);
}