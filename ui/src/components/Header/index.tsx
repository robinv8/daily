import { useState } from 'react'

export default () => {

  const [state, setState] = useState(false)

  // Replace javascript:void(0) path with your path
  const navigation = [
      { title: "Customers", path: "javascript:void(0)" },
      { title: "Careers", path: "javascript:void(0)" },
      { title: "Guides", path: "javascript:void(0)" },
      { title: "Partners", path: "javascript:void(0)" }
  ]

  return (
      <nav className="bg-white w-full border-b md:border-0 md:static">
          <div className="items-center px-4 max-w-screen-2xl mx-auto md:flex md:px-8">
              <div className="flex items-center justify-between py-3 md:py-5 md:block">
                    <h1 className='text-3xl'>
                        日常精选
                    </h1>
              </div>
           
          </div>
      </nav>
  )
}