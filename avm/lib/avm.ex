defmodule Avm do
  # Pin 2 works on any pico device with an external LED
  # @pin 2
  # Comment out above and uncomment below to use Pico W onboard LED
  @pin {:wl, 0}

  def start() do
    loop(@pin, :low)
  end

  defp loop(pin, level) do
    # :io.format(~c"Setting pin ~p ~p~n", [pin, level])
    GPIO.digital_write(pin, level)
    Process.sleep(200)
    loop(pin, toggle(level))
  end

  defp toggle(:high) do
    :low
  end

  defp toggle(:low) do
    :high
  end
end
